// Package rehapt allows building REST HTTP API test cases by describing the request to execute
// and the expected response body. The library takes care of comparing the expected and actual response
// and reports any errors.
// It has been designed to work very well for JSON APIs
//
// Example:
//
//	func TestAPISimple(t *testing.T) {
//	  r := NewRehapt(t, yourHttpServerMux)
//
//	  // Each testcase consist of a description of the request to execute
//	  // and a description of the expected response
//	  // By default the response description is exhaustive.
//	  // If an actual response field is not listed here, an error will be triggered
//	  // of course if an expected field described here is not present in response, an error will be triggered too.
//	  r.TestAssert(TestCase{
//	      Request: TestRequest{
//	          Method: "GET",
//	          Path:   "/api/user/1",
//	      },
//	      Response: TestResponse{
//	          Code: http.StatusOK,
//	          Body: M{
//	              "id":   "1",
//	              "name": "John",
//	              "age":  51,
//	              "pets": S{ // S for slice, M for map. Easy right ?
//	                  M{
//	                      "id":   "2",
//	                      "name": "Pepper the cat",
//	                      "type": "cat",
//	                  },
//	              },
//	              "weddingdate": "2019-06-22T16:00:00.000Z",
//	          },
//	      },
//	  })
//	}
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
	unmarshaler            UnmarshalFn
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

// NewRehapt builds a new Rehapt instance from the given http.Handler.
// `handler` must be your server global handler. For example, it could be
// a simple http.NewServeMux() or a complex third-party library mux.
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

// SetHttpHandler allows changing the http.Handler used to run requests
func (r *Rehapt) SetHttpHandler(handler http.Handler) {
	r.httpHandler = handler
}

// SetMarshaler allows changing the marshal function used to encode the request body.
// The default marshaler is json.Marshal
func (r *Rehapt) SetMarshaler(marshaler func(v interface{}) ([]byte, error)) {
	r.marshaler = marshaler
}

// SetUnmarshaler allows changing the unmarshal function used to decode the response.
// The default unmarshaler is json.Unmarshal
func (r *Rehapt) SetUnmarshaler(unmarshaler func(data []byte, v interface{}) error) {
	r.unmarshaler = unmarshaler
}

// SetErrorHandler allows changing the object handling errors which is called when TestAssert() encounters an error.
// Setting ErrorHandler to nil will simply print the errors on stdout
func (r *Rehapt) SetErrorHandler(errorHandler ErrorHandler) {
	r.errorHandler = errorHandler
}

// SetDefaultHeaders allows setting all default request headers.
// These headers will be added to all requests, however each
// TestCase can override their values
func (r *Rehapt) SetDefaultHeaders(headers http.Header) {
	// nil defaultHeaders could lead to panic later
	if headers != nil {
		r.defaultHeaders = headers
	}
}

// GetDefaultHeaders allows getting all default request headers.
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

// SetDefaultHeader allows setting a default request header.
// This header will be added to all requests, however each
// TestCase can override its value
func (r *Rehapt) SetDefaultHeader(name string, value string) {
	r.defaultHeaders.Set(name, value)
}

// AddDefaultHeader allows adding a default request header.
// This header will be added to all requests, however each
// TestCase can override its value
func (r *Rehapt) AddDefaultHeader(name string, value string) {
	r.defaultHeaders.Add(name, value)
}

// SetDefaultTimeDeltaFormat allows changing the default time format
// It is used by TimeDelta, to parse the actual string value as a time.Time
// Default is set to time.RFC3339 which is ok for JSON.
// This default format can be changed manually for each TimeDelta
func (r *Rehapt) SetDefaultTimeDeltaFormat(format string) {
	r.defaultTimeDeltaFormat = format
}

// Test is the main function of the library
// it executes a given TestCase, i.e. does the request and
// checks if the actual response matches the expected response
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
		marshaler := r.marshaler
		if testcase.Request.BodyMarshaler != nil {
			marshaler = testcase.Request.BodyMarshaler
		}

		bodyData, err := marshaler(testcase.Request.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal the testcase request body. %v", err)
		}
		body = bytes.NewBuffer(bodyData)
	}

	// Path should be either a string or a ReplaceFn
	requestPath := ""
	if repl, ok := testcase.Request.Path.(ReplaceFn); ok == true {
		requestPath, err = repl(r)
		if err != nil {
			return fmt.Errorf("failed to replace path. %v", err)
		}
	} else if p, ok := testcase.Request.Path.(string); ok == true {
		// Default to auto-replace
		requestPath, err = r.replaceVars(p)
		if err != nil {
			return fmt.Errorf("error while replacing variables in path. %v", err)
		}
	} else {
		return fmt.Errorf("invalid path type %T, only string or rehapt.ReplaceFn supported", testcase.Request.Path)
	}

	// Now start to build the HTTP request
	request, err := http.NewRequest(testcase.Request.Method, requestPath, body)
	if err != nil {
		return fmt.Errorf("failed to build HTTP request. %v", err)
	}

	// Add the default headers (if any)
	request.Header = cloneHeader(r.defaultHeaders)

	// Add the testcase defined headers. This overrides any default header previously set
	for k, values := range testcase.Request.Headers {
		request.Header.Del(k)
		for _, value := range values {
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
	if err := r.compare(testcase.Response.Code, response.StatusCode); err != nil {
		codeError = fmt.Errorf("response code does not match. %v", err)
	}

	// Check headers if requested
	if testcase.Response.Headers != nil {
		if err := r.compare(testcase.Response.Headers, response.Header); err != nil {
			headersError = fmt.Errorf("response headers do not match. %v", err)
		}
	}

	bodyError = func() error {
		var responseBody interface{}
		if response.Body != nil {
			data, err := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				return fmt.Errorf("cannot read response body. %v", err)
			}

			if len(data) > 0 {
				unmarshaler := r.unmarshaler
				if testcase.Response.BodyUnmarshaler != nil {
					unmarshaler = testcase.Response.BodyUnmarshaler
				}

				if err := unmarshaler(data, &responseBody); err != nil {
					// If body is nil, then continue with nil decoded body
					// the compare function will handle if that's expected or not
					// but we don't want to report an unmarshal error
					if err != io.EOF {
						return fmt.Errorf("cannot unmarshal response body. %v", err)
					}
				}
			}
		}

		// Compare the response body with our testcase response body
		// We could have used reflect.DeepEqual but we want finer comparison,
		// which allow ignoring some fields, storing variables, using variables, etc.
		// This is the main purpose of this library
		if err := r.compare(testcase.Response.Body, responseBody); err != nil {
			return err
		}

		return nil
	}()

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
				// End of call stack
				break
			}

			// retrieve function name from prog-counter
			function := runtime.FuncForPC(pc)
			if function == nil {
				break
			}

			// functionName will have form package.FuncName
			// "github.com/thib-ack/rehapt_test.TestErrStringResponseBody"
			functionName := function.Name()

			// That's the std testing library
			// which is calling the tests
			if functionName == "testing.tRunner" {
				// Normally we break here, when we reach the testing lib
				break
			}

			filename := path.Base(file)
			callingStack = append(callingStack, fmt.Sprintf("%v:%d: %v", filename, line, functionName))
		}

		message := fmt.Sprintf("%v\nError: %v", strings.Join(callingStack, "\n"), err)

		if r.errorHandler != nil {
			// Start with a \n because testing.T Errorf() prints data and does not start on a new line
			r.errorHandler.Errorf("\n%s", message)
		} else {
			fmt.Printf("%s\n", message)
		}
	}
}

func (r *Rehapt) initComparators() {
	// Fill the list of supported comparators
	// Note the list order does matter because
	// first matching comparator is used.
	r.comparators = []comparator{
		{
			ExpectedKind: reflect.Struct,
			ExpectedType: reflect.TypeOf(time.Time{}),
			Compare:      r.timeCompare,
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
	// shortcut: nil == nil
	if expected == nil && actual == nil {
		return nil
	}
	// shortcut too
	if expected == nil {
		return fmt.Errorf("expected is nil but got %v", actual)
	}

	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	ctx := compareCtx{
		Expected:      expected,
		ExpectedKind:  expectedType.Kind(),
		ExpectedType:  expectedType,
		ExpectedValue: reflect.ValueOf(expected),
		Actual:        actual,
		ActualType:    actualType,
		ActualValue:   reflect.ValueOf(actual),
	}

	// If expected is a CompareFn function, then call it
	if cmp, ok := expected.(CompareFn); ok == true {
		return cmp(r, ctx)
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
