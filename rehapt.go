// Package rehapt allows to build REST HTTP API test cases by describing the request to execute
// and the expected response object. The library takes care of comparing the expected and actual response
// and reports any errors.
// It has been designed to work very well for JSON APIs
//
// Example:
//
//  func TestAPISimple(t *testing.T) {
//    r := NewRehapt(yourHttpServerMux)
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
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"time"
)

// DefaultFailFunction is the default fonction called by
// TestAssert in case of failure. It simply fmt.Println() the error.
func DefaultFailFunction(err error) {
	fmt.Println("Error:", err)
}

// Rehapt - REST HTTP API Test
//
// This is the main structure of the library.
// You can build it using the NewRehapt() function.
type Rehapt struct {
	httpHandler            http.Handler
	marshaler              func(v interface{}) ([]byte, error)
	unmarshaler            func(data []byte, v interface{}) error
	fail                   func(err error)
	defaultHeaders         map[string]string
	variables              map[string]interface{}
	defaultTimeDeltaFormat string
	variableStoreRegexp    *regexp.Regexp
	variableLoadRegexp     *regexp.Regexp
	variableNameRegexp     *regexp.Regexp
}

// NewRehapt build a new Rehapt instance from the given http.Handler.
// `handler` must be your server global handler. For example it could be
// a simple http.NewServeMux() or an complex third-party library mux
func NewRehapt(handler http.Handler) *Rehapt {
	return &Rehapt{
		httpHandler:            handler,
		marshaler:              json.Marshal,
		unmarshaler:            json.Unmarshal,
		fail:                   DefaultFailFunction,
		defaultHeaders:         make(map[string]string),
		variables:              make(map[string]interface{}),
		defaultTimeDeltaFormat: time.RFC3339,
		variableStoreRegexp:    regexp.MustCompile(`^\$([a-zA-Z0-9]+)\$$`),
		variableLoadRegexp:     regexp.MustCompile(`_[a-zA-Z0-9]+_`),
		variableNameRegexp:     regexp.MustCompile(`^[a-zA-Z0-9]+$`),
	}
}

// SetHttpHandler allow to change the http.Handler used to run requests
func (r *Rehapt) SetHttpHandler(handler http.Handler) {
	r.httpHandler = handler
}

// SetMarshaler allow to change the marshaling function used to encode requests body.
// The default marshaler is json.Marshal
func (r *Rehapt) SetMarshaler(marshaler func(v interface{}) ([]byte, error)) {
	r.marshaler = marshaler
}

// SetUnmarshaler allow to change the unmarshaling function used to decode requests response.
// The default unmarshaler is json.Unmarshal
func (r *Rehapt) SetUnmarshaler(unmarshaler func(data []byte, v interface{}) error) {
	r.unmarshaler = unmarshaler
}

// SetFail allow to change the function called when TestAssert() encounter an error.
// The default Fail callback is DefaultFailFunction which simply prints the error
func (r *Rehapt) SetFail(fail func(err error)) {
	r.fail = fail
}

// GetVariable allow to retrive a variable value from its name.
// nil is returned if variable is not found
func (r *Rehapt) GetVariable(name string) interface{} {
	return r.variables[name]
}

// GetVariableString allow to retrive a variable value as a string from its name
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

// GetDefaultHeader returns the default request header value from its name.
// Default headers are added automatically to all requests
func (r *Rehapt) GetDefaultHeader(name string) string {
	return r.defaultHeaders[name]
}

// SetDefaultHeader allow to add a default request header.
// This header will be added to all requests, however each
// TestCase can override its value
func (r *Rehapt) SetDefaultHeader(name string, value string) {
	r.defaultHeaders[name] = value
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
	// If a body has been defined, then marshal it
	if testcase.Request.Body != nil {
		bodyData, err := r.marshaler(testcase.Request.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal the testcase request body. %v", err)
		}
		body = bytes.NewBuffer(bodyData)

	} else if testcase.Request.Raw != nil {
		// If a raw body has been defined use it as-is
		body = testcase.Request.Raw
	}

	// The path might contains a variable reference (like _xx_). we have to replace it.
	testcase.Request.Path = r.replaceVars(testcase.Request.Path)

	// Now start to build the HTTP request
	request, err := http.NewRequest(testcase.Request.Method, testcase.Request.Path, body)
	if err != nil {
		return fmt.Errorf("failed to build HTTP request. %v", err)
	}

	// Add the default headers (if any)
	for k, v := range r.defaultHeaders {
		request.Header.Set(k, v)
	}
	// Add the testcase defined headers. This overrides any default header previously set
	for k, v := range testcase.Request.Headers {
		request.Header.Set(k, v)
	}

	// Now execute the request and record its response
	recorder := httptest.NewRecorder()
	r.httpHandler.ServeHTTP(recorder, request)

	// And start to check result.

	// First check HTTP response code
	// Maybe we have to ignore this completely as requested by the user
	if testcase.Response.Code != AnyCode {
		if testcase.Response.Code != recorder.Code {
			return fmt.Errorf("response code does not match. Expected %d, got %d", testcase.Response.Code, recorder.Code)
		}
	}

	// Check headers, but not all of them. Only the ones expected by the user
	for header, value := range testcase.Response.Headers {
		actualValue := recorder.Header().Get(header)
		if value != actualValue {
			return fmt.Errorf("response header %v does not match. Expected %v, got %v", header, value, actualValue)
		}
	}

	// Want a raw comparison ?
	// This is useful if response cannot be unmarshaled. (for example simple plain/text output)
	if testcase.Response.Raw != nil {
		actualBody := recorder.Body.String()

		// However we still provide some features for this.
		// We can use Regexp, RegexpVars or simple string
		switch rawObject := testcase.Response.Raw.(type) {
		case Regexp:
			re, err := regexp.Compile(string(rawObject))
			if err != nil {
				return err
			}
			if re.MatchString(actualBody) == false {
				return fmt.Errorf("regexp '%v' does not match '%v'", rawObject, actualBody)
			}
			return nil

		case RegexpVars:
			re, err := regexp.Compile(rawObject.Regexp)
			if err != nil {
				return err
			}
			elements := re.FindStringSubmatch(actualBody)
			if len(elements) == 0 {
				return fmt.Errorf("regexp '%v' does not match '%v'", rawObject.Regexp, actualBody)
			}

			for groupid, varname := range rawObject.Vars {
				if groupid >= len(elements) {
					return fmt.Errorf("expected variable index %d overflow regexp group count of %d", groupid, len(elements))
				}
				if r.validVarname(varname) == false {
					return fmt.Errorf("invalid variable name %v for group %d", varname, groupid)
				}
				r.storeIfVariable("$"+varname+"$", elements[groupid])
			}
			return nil

		case string:
			if r.storeIfVariable(rawObject, recorder.Body.String()) == true {
				// This was a variable store operation. no comparison to do
				return nil
			}
			// Otherwise, compare full values
			if rawObject != recorder.Body.String() {
				return fmt.Errorf("response body does not match. Expected %v, got %v", rawObject, recorder.Body.String())
			}
			return nil
		}

		return fmt.Errorf("unsupported Raw object type %T", testcase.Response.Raw)
	}

	data, err := ioutil.ReadAll(recorder.Body)
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
}

// TestAssert works exactly like Test except it reports the error if not nil
// using the fail callback function defined by SetFail (or default one)
func (r *Rehapt) TestAssert(testcase TestCase) {
	if err := r.Test(testcase); err != nil {
		if r.fail != nil {
			r.fail(err)
		}
	}
}

func (r *Rehapt) validVarname(name string) bool {
	return r.variableNameRegexp.MatchString(name)
}

func (r *Rehapt) replaceVars(str string) string {
	return r.variableLoadRegexp.ReplaceAllStringFunc(str, func(name string) string {
		// Remove the '_' prefix and suffix to get only the variable name
		varname := name[1 : len(name)-1]
		if value, ok := r.variables[varname]; ok == true {
			if str, ok := value.(string); ok == true {
				return str
			}
		}
		return name
	})
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

	// Maybe we have to ignore this completely as requested by the user
	if expectedType == reflect.TypeOf(Any) {
		return nil
	}

	expectedKind := expectedType.Kind()
	actualKind := actualType.Kind()

	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	// We'll have only Slice and Map, never Struct as we unmarshal in interface{}
	switch expectedKind {
	case reflect.Struct:
		// Some of our custom struct

		// Time with delta comparison.
		// actual value must be a string (time like "2012-04-23T18:25:43.511Z")
		if timeDelta, ok := expected.(TimeDelta); ok == true {
			if actualKind != reflect.String {
				return fmt.Errorf("different kinds. Expected string, got %v", actualKind)
			}

			// Use specific time format or default one if not specified
			format := r.defaultTimeDeltaFormat
			if timeDelta.Format != "" {
				format = timeDelta.Format
			}

			actual, err := time.Parse(format, actualValue.String())
			if err != nil {
				return fmt.Errorf("invalid time. %v", err)
			}

			dt := timeDelta.Time.Sub(actual)
			if dt < -timeDelta.Delta || dt > timeDelta.Delta {
				return fmt.Errorf("max difference between %v and %v allowed is %v, but difference was %v", timeDelta.Time, actual, timeDelta.Delta, dt)
			}
			return nil
		}

		if numDelta, ok := expected.(NumberDelta); ok == true {

			actualFloatValue := float64(0.0)
			switch actualKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				actualFloatValue = float64(actualValue.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				actualFloatValue = float64(actualValue.Uint())
			case reflect.Float32, reflect.Float64:
				actualFloatValue = actualValue.Float()
			default:
				return fmt.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", actualKind)
			}

			delta := math.Abs(numDelta.Value - actualFloatValue)
			if delta > numDelta.Delta {
				return fmt.Errorf("max difference between %v and %v allowed is %v, but difference was %v", numDelta.Value, actual, numDelta.Delta, delta)
			}
			return nil
		}

		if reVars, ok := expected.(RegexpVars); ok == true {
			if actualKind != reflect.String {
				return fmt.Errorf("different kinds. Expected string, got %v", actualKind)
			}

			actualStr := actualValue.String()

			re, err := regexp.Compile(reVars.Regexp)
			if err != nil {
				return err
			}
			elements := re.FindStringSubmatch(actualStr)
			if len(elements) == 0 {
				return fmt.Errorf("regexp '%v' does not match '%v'", reVars.Regexp, actualStr)
			}

			for groupid, varname := range reVars.Vars {
				if groupid >= len(elements) {
					return fmt.Errorf("expected variable index %d overflow regexp group count of %d", groupid, len(elements))
				}
				if err := r.SetVariable(varname, elements[groupid]); err != nil {
					return err
				}
			}
			return nil
		}

		return fmt.Errorf("unexpected struct type %v", expectedType)

	case reflect.Slice:
		if actualKind != reflect.Slice {
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		// In case of slice, we have to compare the 2 slices.
		if expectedValue.Len() != actualValue.Len() {
			return fmt.Errorf("different slice sizes. Expected %v, got %v. Expected %v got %v", expectedValue.Len(), actualValue.Len(), expected, actual)
		}

		if expectedType == reflect.TypeOf(UnsortedS{}) {
			// Unordered comparison
			// We build a list of all the indexes (0,1,2,..,N-1)
			// So each time we find a matching element, we can remove its index from this list
			// and ignore it on next search
			actualIndexes := make([]int, actualValue.Len())
			for i := range actualIndexes {
				actualIndexes[i] = i
			}

		nextExpected:
			for i := 0; i < expectedValue.Len(); i++ {
				expectedElement := expectedValue.Index(i)

				// Now find a matching element in actual object.
				// Once found, ignore the index.
				for j := 0; j < len(actualIndexes); j++ {
					idx := actualIndexes[j]
					actualElement := actualValue.Index(idx)

					if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err == nil {
						// Thats a match, ignore this index now, and continue to next expected.
						actualIndexes = append(actualIndexes[:j], actualIndexes[j+1:]...)
						continue nextExpected
					}
				}

				// If we arrive here, we have an expected not matching any actual
				return fmt.Errorf("expected element %v at index %v not found", expectedElement, i)
			}

			// If here we still have actual index, it means unmatched element thats bad
			if len(actualIndexes) > 0 {
				return fmt.Errorf("actual elements at indexes %v not found", actualIndexes)
			}

		} else {
			// ordered comparison
			for i := 0; i < expectedValue.Len(); i++ {
				expectedElement := expectedValue.Index(i)
				actualElement := actualValue.Index(i)
				if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
					return fmt.Errorf("slice element %v does not match. %v", i, err)
				}
			}
		}

		return nil

	case reflect.Map:
		if actualKind != reflect.Map {
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		// In case of map, we have to compare the 2 maps.
		if expectedType.Key() != actualType.Key() {
			return fmt.Errorf("different map key types. Expected %v, got %v", expectedType.Key(), actualType.Key())
		}

		// Partial match. Ignore the keys not listed in expected map
		// to do this we just have to skip the map size comparison
		if expectedType != reflect.TypeOf(PartialM{}) {
			if expectedValue.Len() != actualValue.Len() {
				return fmt.Errorf("different map sizes. Expected %v, got %v", expectedValue.Len(), actualValue.Len())
			}
		}

		keys := expectedValue.MapKeys()
		for _, key := range keys {
			expectedElement := expectedValue.MapIndex(key)
			actualElement := actualValue.MapIndex(key)

			if actualElement.IsValid() == false {
				return fmt.Errorf("expected key %v not found", key)
			}

			if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
				return fmt.Errorf("map element [%v] does not match. %v", key, err)
			}
		}

		return nil

	case reflect.String:
		expectedStr := expectedValue.String()

		if expectedType == reflect.TypeOf(StoreVar("")) {
			// Don't compare but store the actual value using the expectedStr as variable name
			if err := r.SetVariable(expectedStr, actual); err != nil {
				return err
			}
			return nil
		}

		if expectedType == reflect.TypeOf(LoadVar("")) {
			// Compare actual with the loaded value which might not be string
			value := r.GetVariable(expectedStr)
			return r.compare(value, actual)
		}

		if r.storeIfVariable(expectedStr, actual) == true {
			// This was a variable store operation. no comparison to do
			return nil
		}

		if actualKind != reflect.String {
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		actualStr := actualValue.String()

		// Make var replacement in case of
		expectedStr = r.replaceVars(expectedStr)

		// If regexp, then process differently
		if expectedType == reflect.TypeOf(Regexp("")) {
			re, err := regexp.Compile(expectedStr)
			if err != nil {
				return err
			}
			if re.MatchString(actualStr) == false {
				return fmt.Errorf("regexp '%v' does not match '%v'", expectedStr, actualStr)
			}
			return nil
		}

		// classic comparison
		if expectedStr != actualStr {
			return fmt.Errorf("strings does not match. Expected '%v', got '%v'", expectedStr, actualStr)
		}

		return nil

	case reflect.Bool:
		if actualKind != reflect.Bool {
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		expectedBool := expectedValue.Bool()
		actualBool := actualValue.Bool()
		// classic comparison
		if expectedBool != actualBool {
			return fmt.Errorf("bools does not match. Expected %v, got %v", expectedBool, actualBool)
		}

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		expectedInt := expectedValue.Int()

		switch actualKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actualInt := actualValue.Int()
			// classic comparison
			if expectedInt != actualInt {
				return fmt.Errorf("integers does not match. Expected %v, got %v", expectedInt, actualInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actualInt := actualValue.Uint()
			// classic comparison
			if uint64(expectedInt) != actualInt {
				return fmt.Errorf("uintegers does not match. Expected %v, got %v", expectedInt, actualInt)
			}
		case reflect.Float32, reflect.Float64:
			actualFloat := actualValue.Float()
			// classic comparison
			if float64(expectedInt) != actualFloat {
				return fmt.Errorf("floats does not match. Expected %v, got %v", expectedInt, actualFloat)
			}
		default:
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		expectedInt := expectedValue.Uint()

		switch actualKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actualInt := actualValue.Int()
			// classic comparison
			if int64(expectedInt) != actualInt {
				return fmt.Errorf("integers does not match. Expected %v, got %v", expectedInt, actualInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actualInt := actualValue.Uint()
			// classic comparison
			if expectedInt != actualInt {
				return fmt.Errorf("uintegers does not match. Expected %v, got %v", expectedInt, actualInt)
			}
		case reflect.Float32, reflect.Float64:
			actualFloat := actualValue.Float()
			// classic comparison
			if float64(expectedInt) != actualFloat {
				return fmt.Errorf("floats does not match. Expected %v, got %v", expectedInt, actualFloat)
			}
		default:
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		return nil

	case reflect.Float32, reflect.Float64:

		expectedFloat := expectedValue.Float()

		switch actualKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actualInt := actualValue.Int()
			// classic comparison
			if int64(expectedFloat) != actualInt {
				return fmt.Errorf("integers does not match. Expected %v, got %v", expectedFloat, actualInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actualInt := actualValue.Uint()
			// classic comparison
			if uint64(expectedFloat) != actualInt {
				return fmt.Errorf("uintegers does not match. Expected %v, got %v", expectedFloat, actualInt)
			}
		case reflect.Float32, reflect.Float64:
			actualFloat := actualValue.Float()
			// classic comparison
			if expectedFloat != actualFloat {
				return fmt.Errorf("floats does not match. Expected %v, got %v", expectedFloat, actualFloat)
			}
		default:
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		return nil
	}

	return fmt.Errorf("unhandled type %T", expected)
}
