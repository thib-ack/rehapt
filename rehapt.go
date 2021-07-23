// Package rehapt allows to build REST HTTP API test cases by describing the request to execute
// and the expected response object. The library takes care of comparing the expected and actual response
// and reports any errors.
// It has been designed to work very well for JSON APIs
//
// Example:
//
//  func TestAPISimple(t *testing.T) {
//    r := NewRehapt(t, yourHttpServerMux)
//
//    // Each testcase consist of a description of the request to execute
//    // and a description of the expected response
//    // By default the response description is exhaustive.
//    // If an actual response field is not listed here, an error will be triggered
//    // of course if an expected field described here is not present in response, an error will be triggered too.
//    r.TestAssert(TestCase{
//        Request: TestRequest{
//            Method: "GET",
//            Path:   "/api/user/1",
//        },
//        Response: TestResponse{
//            Code: http.StatusOK,
//            Object: M{
//                "id":   "1",
//                "name": "John",
//                "age":  51,
//                "pets": S{ // S for slice, M for map. Easy right ?
//                    M{
//                        "id":   "2",
//                        "name": "Pepper the cat",
//                        "type": "cat",
//                    },
//                },
//                "weddingdate": "2019-06-22T16:00:00.000Z",
//            },
//        },
//    })
//  }
//
// See https://github.com/thib-ack/rehapt/tree/master/examples for more examples
package rehapt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Rehapt - REST HTTP API Test
//
// This is the main structure of the library.
// You can build it using the NewRehapt() function.
type Rehapt struct {
	httpHandler            http.Handler
	marshaler              func(v interface{}) ([]byte, error)
	unmarshaler            func(data []byte, v interface{}) error
	errorHandler           ErrorHandler
	defaultHeaders         http.Header
	variables              map[string]interface{}
	defaultTimeDeltaFormat string
	variableStoreRegexp    *regexp.Regexp
	variableLoadRegexp     *regexp.Regexp
	variableNameRegexp     *regexp.Regexp
	floatPrecision         int
	comparators            []comparator
}

// NewRehapt build a new Rehapt instance from the given http.Handler.
// `handler` must be your server global handler. For example it could be
// a simple http.NewServeMux() or an complex third-party library mux.
// `errorHandler` can be the *testing.T parameter of your test,
// if value is nil, the errors are printed on stdout
func NewRehapt(errorHandler ErrorHandler, handler http.Handler) *Rehapt {
	r := &Rehapt{
		httpHandler:            handler,
		marshaler:              json.Marshal,
		unmarshaler:            json.Unmarshal,
		errorHandler:           errorHandler,
		defaultHeaders:         make(http.Header),
		variables:              make(map[string]interface{}),
		defaultTimeDeltaFormat: time.RFC3339,
		variableStoreRegexp:    regexp.MustCompile(`^\$([a-zA-Z0-9]+)\$$`),
		variableLoadRegexp:     regexp.MustCompile(`_([a-zA-Z0-9]+)_`),
		variableNameRegexp:     regexp.MustCompile(`^[a-zA-Z0-9]+$`),
		floatPrecision:         -1,
		comparators:            nil,
	}
	r.initComparators()
	return r
}

// SetHttpHandler allow to change the http.Handler used to run requests
func (r *Rehapt) SetHttpHandler(handler http.Handler) {
	r.httpHandler = handler
}

// SetMarshaler allow to change the marshal function used to encode requests body.
// The default marshaler is json.Marshal
func (r *Rehapt) SetMarshaler(marshaler func(v interface{}) ([]byte, error)) {
	r.marshaler = marshaler
}

// SetUnmarshaler allow to change the unmarshal function used to decode requests response.
// The default unmarshaler is json.Unmarshal
func (r *Rehapt) SetUnmarshaler(unmarshaler func(data []byte, v interface{}) error) {
	r.unmarshaler = unmarshaler
}

// SetErrorHandler allow to change the object handling errors which is called when TestAssert() encounter an error.
// Setting ErrorHandler to nil will simply print the errors on stdout
func (r *Rehapt) SetErrorHandler(errorHandler ErrorHandler) {
	r.errorHandler = errorHandler
}

// GetVariable allow to retrieve a variable value from its name.
// nil is returned if variable is not found
func (r *Rehapt) GetVariable(name string) interface{} {
	return r.variables[name]
}

// GetVariableString allow to retrieve a variable value as a string from its name
// empty string is returned if variable is not found
func (r *Rehapt) GetVariableString(name string) string {
	if value, ok := r.variables[name].(string); ok == true {
		return value
	}
	return ""
}

// SetVariable allow to define manually a variable.
// Variable names are strings, however values can be any type
func (r *Rehapt) SetVariable(name string, value interface{}) error {
	if r.validVarname(name) == false {
		return fmt.Errorf("invalid variable name %v", name)
	}
	r.variables[name] = value
	return nil
}

// SetDefaultHeaders allow to set all default request headers.
// These headers will be added to all requests, however each
// TestCase can override their values
func (r *Rehapt) SetDefaultHeaders(headers http.Header) {
	r.defaultHeaders = headers
}

// GetDefaultHeaders allow to get all default request headers.
// These headers will be added to all requests, however each
// TestCase can override their values
func (r *Rehapt) GetDefaultHeaders() http.Header {
	return r.defaultHeaders
}

// GetDefaultHeader returns the default request header value from its name.
// Default headers are added automatically to all requests
func (r *Rehapt) GetDefaultHeader(name string) string {
	return r.defaultHeaders.Get(name)
}

// SetDefaultHeader allow to set a default request header.
// This header will be added to all requests, however each
// TestCase can override its value
func (r *Rehapt) SetDefaultHeader(name string, value string) {
	r.defaultHeaders.Set(name, value)
}

// AddDefaultHeader allow to add a default request header.
// This header will be added to all requests, however each
// TestCase can override its value
func (r *Rehapt) AddDefaultHeader(name string, value string) {
	r.defaultHeaders.Add(name, value)
}

// SetDefaultTimeDeltaFormat allow to change the default time format
// It is used by TimeDelta, to parse the actual string value as a time.Time
// Default is set to time.RFC3339 which is ok for JSON.
// This default format can be changed manually for each TimeDelta
func (r *Rehapt) SetDefaultTimeDeltaFormat(format string) {
	r.defaultTimeDeltaFormat = format
}

// SetStoreShortcutBounds modify the strings used as prefix and suffix to identify
// a shortcut version of the store variable operation. The default prefix and suffix is "$" which makes
// the default shortcut form like "$myvar$".
func (r *Rehapt) SetStoreShortcutBounds(prefix string, suffix string) error {
	if prefix == "" {
		return fmt.Errorf("invalid prefix, cannot be empty")
	}
	if suffix == "" {
		return fmt.Errorf("invalid suffix, cannot be empty")
	}
	prefixEscaped := regexp.QuoteMeta(prefix)
	suffixEscaped := regexp.QuoteMeta(suffix)
	re, err := regexp.Compile(`^` + prefixEscaped + `([a-zA-Z0-9]+)` + suffixEscaped + `$`)
	if err != nil {
		return err
	}
	r.variableStoreRegexp = re
	return nil
}

// SetLoadShortcutBounds modify the strings used as prefix and suffix to identify
// a shortcut version of the load variable operation. The default prefix and suffix is "_" which makes
// the default shortcut form like "_myvar_".
func (r *Rehapt) SetLoadShortcutBounds(prefix string, suffix string) error {
	if prefix == "" {
		return fmt.Errorf("invalid prefix, cannot be empty")
	}
	if suffix == "" {
		return fmt.Errorf("invalid suffix, cannot be empty")
	}
	prefixEscaped := regexp.QuoteMeta(prefix)
	suffixEscaped := regexp.QuoteMeta(suffix)
	re, err := regexp.Compile(prefixEscaped + `([a-zA-Z0-9]+)` + suffixEscaped)
	if err != nil {
		return err
	}
	r.variableLoadRegexp = re
	return nil
}

// SetLoadShortcutFloatPrecision change the precision of float formatting when
// used with a load shortcut. For example "value is _myvar_" can be replaced by
// "value is 10.50" or "value is 10.500000".
func (r *Rehapt) SetLoadShortcutFloatPrecision(precision int) {
	r.floatPrecision = precision
}

// Test is the main function of the library
// it executes a given TestCase, i.e. do the request and
// check if the actual response is matching the expected response
func (r *Rehapt) Test(testcase TestCase) error {
	// If we don't have the minimum, we cannot go further.
	if r.httpHandler == nil {
		return fmt.Errorf("nil HTTP handler")
	}
	if r.marshaler == nil {
		return fmt.Errorf("nil marshaler")
	}
	if r.unmarshaler == nil {
		return fmt.Errorf("nil unmarshaler")
	}
	if testcase.Request.Method == "" {
		return fmt.Errorf("incomplete testcase. Missing HTTP method")
	}
	if testcase.Request.Path == "" {
		return fmt.Errorf("incomplete testcase. Missing URL path")
	}

	var body io.Reader
	var err error
	// If a body has been defined, then marshal it
	if testcase.Request.Body != nil {
		bodyData, err := r.marshaler(testcase.Request.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal the testcase request body. %v", err)
		}
		body = bytes.NewBuffer(bodyData)

	} else if testcase.Request.RawBody != nil {
		// If a raw body has been defined use it as-is (no marshal operation)
		// unless variable replacement is allowed
		if testcase.Request.NoRawBodyVariableReplacement == true {
			body = testcase.Request.RawBody
		} else {
			// This could be optimized
			rawBody, err := ioutil.ReadAll(testcase.Request.RawBody)
			if err != nil {
				return fmt.Errorf("error while reading raw body. %v", err)
			}
			rawBodyStr, err := r.replaceVars(string(rawBody))
			if err != nil {
				return fmt.Errorf("error while replacing variables in raw body. %v", err)
			}
			body = bytes.NewBufferString(rawBodyStr)
		}
	}

	// The path might contains a variable reference (like _xx_). we have to replace it.
	if testcase.Request.NoPathVariableReplacement == false {
		testcase.Request.Path, err = r.replaceVars(testcase.Request.Path)
		if err != nil {
			return fmt.Errorf("error while replacing variables in path. %v", err)
		}
	}

	// Now start to build the HTTP request
	request, err := http.NewRequest(testcase.Request.Method, testcase.Request.Path, body)
	if err != nil {
		return fmt.Errorf("failed to build HTTP request. %v", err)
	}

	// Add the default headers (if any)
	request.Header = cloneHeader(r.defaultHeaders)

	// Add the testcase defined headers. This overrides any default header previously set
	for k, values := range testcase.Request.Headers {
		if testcase.Request.NoHeadersVariableReplacement == false {
			k, err = r.replaceVars(k)
			if err != nil {
				return fmt.Errorf("error while replacing variables in header name. %v", err)
			}
		}
		request.Header.Del(k)
		for _, value := range values {
			if testcase.Request.NoHeadersVariableReplacement == false {
				value, err = r.replaceVars(value)
				if err != nil {
					return fmt.Errorf("error while replacing variables in header value. %v", err)
				}
			}
			request.Header.Add(k, value)
		}
	}

	// Now execute the request and record its response
	recorder := httptest.NewRecorder()
	r.httpHandler.ServeHTTP(recorder, request)
	response := recorder.Result()

	// And start to check result.
	// But don't stop on first error, for example if http code doesn't match,
	// we can still compare headers and body.
	var codeError error
	var headersError error
	var bodyError error

	// First check HTTP response code
	// Maybe we have to ignore this completely as requested by the user
	if testcase.Response.Code != AnyCode {
		if testcase.Response.Code != response.StatusCode {
			codeError = fmt.Errorf("response code does not match. Expected %d, got %d", testcase.Response.Code, response.StatusCode)
		}
	}

	// Check headers if requested
	if testcase.Response.Headers != nil {
		if err := r.compare(testcase.Response.Headers, response.Header); err != nil {
			headersError = fmt.Errorf("response headers does not match. %v", err)
		}
	}

	// Want a raw comparison ?
	// This is useful if response cannot be unmarshal. (for example simple plain/text output)
	if testcase.Response.RawBody != nil {
		if err := r.compare(testcase.Response.RawBody, recorder.Body.String()); err != nil {
			bodyError = err
		}

	} else {
		// Use Object expected body
		bodyError = func() error {
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("cannot read response body. %v", err)
			}

			var responseBody interface{}
			if len(data) > 0 {
				if err := r.unmarshaler(data, &responseBody); err != nil {
					// If body is nil, then continue with nil decoded body
					// the compare function will handle if that's expected or not
					// but we don't want to report an unmarshal error
					if err != io.EOF {
						return fmt.Errorf("cannot unmarshal response body. %v", err)
					}
				}
			}

			// Compare the response body with our testcase response object
			// We could have used reflect.DeepEqual but we want finer comparison,
			// which allow ignoring some fields, storing variables, using variables, etc.
			// This is the main purpose of this library
			if err := r.compare(testcase.Response.Object, responseBody); err != nil {
				return err
			}

			return nil
		}()
	}

	// Build an error based on the 3 possible errors on code, headers and body
	if codeError != nil || headersError != nil || bodyError != nil {
		e := ""
		if codeError != nil {
			e += codeError.Error() + "\n"
		}
		if headersError != nil {
			e += headersError.Error() + "\n"
		}
		if bodyError != nil {
			e += bodyError.Error()
		}
		return errors.New(strings.TrimSuffix(e, "\n"))
	}
	return nil
}

// TestAssert works exactly like Test except it reports the error if not nil
// using the ErrorHandler Errorf() function
func (r *Rehapt) TestAssert(testcase TestCase) {
	if err := r.Test(testcase); err != nil {
		// index 0 is this function calling runtime.Caller() -> we can skip it
		// start at index 1 to get the user function calling rehapt.TestAssert()
		//
		// We could use only the index 2, but if somebody is using rehapt.TestAssert() inside another function
		// then it is still good to go further and return all callers recursively until we reach the std testing library
		var callingStack []string
		for i := 1; i < 20; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				// End of call-stack
				break
			}

			// retrieve function name from prog-counter
			function := runtime.FuncForPC(pc)
			if function == nil {
				break
			}

			// functionName will have form package.FuncName
			// "github.com/thib-ack/rehapt_test.TestErrStringResponseObject"
			functionName := function.Name()

			// That's the std testing library
			// which is calling the tests
			if functionName == "testing.tRunner" {
				// Normally we break here, when we reached the testing lib
				break
			}

			filename := path.Base(file)
			callingStack = append(callingStack, fmt.Sprintf("%v:%d: %v", filename, line, functionName))
		}

		message := fmt.Sprintf("%v\nError: %v", strings.Join(callingStack, "\n"), err)

		if r.errorHandler != nil {
			// Start with a \n because testing.T Errorf() prints data and do not start on new line
			r.errorHandler.Errorf("\n" + message)
		} else {
			fmt.Printf(message + "\n")
		}
	}
}

func (r *Rehapt) validVarname(name string) bool {
	return r.variableNameRegexp.MatchString(name)
}

func (r *Rehapt) replaceVars(str string) (string, error) {
	matches := r.variableLoadRegexp.FindAllStringSubmatchIndex(str, -1)
	if len(matches) == 0 {
		return str, nil
	}

	replaced := make([]byte, 0, len(str)*2)
	offset := 0
	for _, match := range matches {
		// Match should be 4 elements
		// For example, "the _var_ move" should return [4, 9, 5, 8] :
		//  0   45  89
		// "the _var_ move"
		// with the 4 indexes : [prefix start, suffix end, varname start, varname end]
		if len(match) < 4 {
			continue
		}
		prefix := match[0]
		suffix := match[1]
		varnameStart := match[2]
		varnameEnd := match[3]

		// remove the prefix and suffix
		varname := str[varnameStart:varnameEnd]
		value := ""

		// Make sure variable exists, or report error
		ivalue, ok := r.variables[varname]
		if ok == false {
			return "", fmt.Errorf("variable %v is not defined", varname)
		}

		// Try to convert value to string
		switch ival := ivalue.(type) {
		case string:
			value = ival
		case int:
			value = strconv.FormatInt(int64(ival), 10)
		case int8:
			value = strconv.FormatInt(int64(ival), 10)
		case int16:
			value = strconv.FormatInt(int64(ival), 10)
		case int32:
			value = strconv.FormatInt(int64(ival), 10)
		case int64:
			value = strconv.FormatInt(ival, 10)
		case uint:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint8:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint16:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint32:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint64:
			value = strconv.FormatUint(ival, 10)
		case float32:
			value = strconv.FormatFloat(float64(ival), 'f', r.floatPrecision, 32)
		case float64:
			value = strconv.FormatFloat(ival, 'f', r.floatPrecision, 64)
		case bool:
			value = strconv.FormatBool(ival)
		default:
			return "", fmt.Errorf("variable %v of type %T cannot be using inside string", varname, ivalue)
		}

		replaced = append(replaced, str[offset:prefix]...)
		replaced = append(replaced, value...)
		offset = suffix
	}

	// Finish with end of str, if any
	if offset < len(str) {
		replaced = append(replaced, str[offset:]...)
	}

	return string(replaced), nil
}

func (r *Rehapt) storeIfVariable(expected string, actual interface{}) bool {
	elements := r.variableStoreRegexp.FindStringSubmatch(expected)
	if len(elements) > 1 {
		// index 0 is the full match.
		// index 1 is the first group, our variable name without the '_' prefix and suffix
		varname := elements[1]
		// We override any stored value
		r.variables[varname] = actual
		return true
	}
	return false
}

func (r *Rehapt) initComparators() {
	// Fill the list of supported comparators
	// Note the list order do matter because
	// first matching comparator is used.
	r.comparators = []comparator{
		{
			ExpectedKind: reflect.Struct,
			ExpectedType: reflect.TypeOf(Not{}),
			Compare:      r.notCompare,
		},
		{
			ExpectedKind: reflect.Struct,
			ExpectedType: reflect.TypeOf(TimeDelta{}),
			Compare:      r.timeDeltaCompare,
		},
		{
			ExpectedKind: reflect.Struct,
			ExpectedType: reflect.TypeOf(NumberDelta{}),
			Compare:      r.numberDeltaCompare,
		},
		{
			ExpectedKind: reflect.Struct,
			ExpectedType: reflect.TypeOf(RegexpVars{}),
			Compare:      r.regexpVarsCompare,
		},
		{
			ExpectedKind: reflect.Slice,
			ExpectedType: reflect.TypeOf(UnsortedS{}),
			Compare:      r.unsortedSliceCompare,
		},
		{
			ExpectedKind: reflect.Slice,
			ExpectedType: nil,
			Compare:      r.sliceCompare,
		},
		{
			ExpectedKind: reflect.Map,
			ExpectedType: reflect.TypeOf(PartialM{}),
			Compare:      r.partialMapCompare,
		},
		{
			ExpectedKind: reflect.Map,
			ExpectedType: nil,
			Compare:      r.mapCompare,
		},
		{
			ExpectedKind: reflect.String,
			ExpectedType: reflect.TypeOf(Any),
			Compare:      r.anyCompare,
		},
		{
			ExpectedKind: reflect.String,
			ExpectedType: reflect.TypeOf(StoreVar("")),
			Compare:      r.storeVarCompare,
		},
		{
			ExpectedKind: reflect.String,
			ExpectedType: reflect.TypeOf(LoadVar("")),
			Compare:      r.loadVarCompare,
		},
		{
			ExpectedKind: reflect.String,
			ExpectedType: reflect.TypeOf(Regexp("")),
			Compare:      r.regexpCompare,
		},
		{
			ExpectedKind: reflect.String,
			ExpectedType: nil,
			Compare:      r.stringCompare,
		},
		{
			ExpectedKind: reflect.Bool,
			ExpectedType: nil,
			Compare:      r.boolCompare,
		},
		{
			ExpectedKind: reflect.Int,
			ExpectedType: nil,
			Compare:      r.intCompare,
		},
		{
			ExpectedKind: reflect.Int8,
			ExpectedType: nil,
			Compare:      r.intCompare,
		},
		{
			ExpectedKind: reflect.Int16,
			ExpectedType: nil,
			Compare:      r.intCompare,
		},
		{
			ExpectedKind: reflect.Int32,
			ExpectedType: nil,
			Compare:      r.intCompare,
		},
		{
			ExpectedKind: reflect.Int64,
			ExpectedType: nil,
			Compare:      r.intCompare,
		},
		{
			ExpectedKind: reflect.Uint,
			ExpectedType: nil,
			Compare:      r.uintCompare,
		},
		{
			ExpectedKind: reflect.Uint8,
			ExpectedType: nil,
			Compare:      r.uintCompare,
		},
		{
			ExpectedKind: reflect.Uint16,
			ExpectedType: nil,
			Compare:      r.uintCompare,
		},
		{
			ExpectedKind: reflect.Uint32,
			ExpectedType: nil,
			Compare:      r.uintCompare,
		},
		{
			ExpectedKind: reflect.Uint64,
			ExpectedType: nil,
			Compare:      r.uintCompare,
		},
		{
			ExpectedKind: reflect.Float32,
			ExpectedType: nil,
			Compare:      r.floatCompare,
		},
		{
			ExpectedKind: reflect.Float64,
			ExpectedType: nil,
			Compare:      r.floatCompare,
		},
	}
}

func (r *Rehapt) compare(expected interface{}, actual interface{}) error {
	// This is perfectly valid
	if expected == nil && actual == nil {
		return nil
	}
	// but this is not. We cannot go further in these 2 cases as there are nothing to compare
	if expected == nil {
		return fmt.Errorf("expected is nil but got %v", actual)
	}
	if actual == nil {
		return fmt.Errorf("expected %v but got nil", expected)
	}

	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	ctx := compareCtx{
		Expected:      expected,
		ExpectedKind:  expectedType.Kind(),
		ExpectedType:  expectedType,
		ExpectedValue: reflect.ValueOf(expected),
		Actual:        actual,
		ActualKind:    actualType.Kind(),
		ActualType:    actualType,
		ActualValue:   reflect.ValueOf(actual),
	}

	// Now find a matching comparator and let it do the job.
	// We iterate through our defined comparators and stop on the first matching one.
	// Either the Kind *and* the Type have to match (for example Kind==String and Type==Regexp)
	// or only the Kind as a generic fallback (for example Kind==String)
	for _, comparator := range r.comparators {
		if comparator.ExpectedKind == ctx.ExpectedKind {
			if comparator.ExpectedType == expectedType || comparator.ExpectedType == nil {
				return comparator.Compare(ctx)
			}
		}
	}
	return fmt.Errorf("unhandled type %T", expected)
}

func cloneHeader(header http.Header) http.Header {
	// Clone() method of http.Header is available only since 1.13
	if header == nil {
		return nil
	}
	clone := make(http.Header, len(header))
	for k, v := range header {
		sl := make([]string, len(v))
		copy(sl, v)
		clone[k] = sl
	}
	return clone
}
