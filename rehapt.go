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

func DefaultFailFunction(err error) {
	fmt.Println("Error:", err)
}

// ideas to add:
// Raw body to request
// replace vars in request body

// Rehapt - REST HTTP API Test
type Rehapt struct {
	httpHandler            http.Handler
	marshaler              func(v interface{}) ([]byte, error)
	unmarshaler            func(data []byte, v interface{}) error
	fail                   func(err error)
	defaultHeaders         map[string]string
	variables              map[string]interface{}
	defaultTimeDeltaFormat string
}

// NewRehapt build a new Rehapt instance from the given http.Handler
// `handler` must be your server global handler. For example it could be
// a simple `http.NewServeMux()` or an complex third-party library mux
func NewRehapt(handler http.Handler) *Rehapt {
	return &Rehapt{
		httpHandler:            handler,
		marshaler:              json.Marshal,
		unmarshaler:            json.Unmarshal,
		fail:                   DefaultFailFunction,
		defaultHeaders:         make(map[string]string),
		variables:              make(map[string]interface{}),
		defaultTimeDeltaFormat: time.RFC3339,
	}
}

// SetHttpHandler allow to change the http.Handler used to run requests
func (r *Rehapt) SetHttpHandler(handler http.Handler) {
	r.httpHandler = handler
}

// SetMarshaler allow to change the marshaling function used to encode requests body
func (r *Rehapt) SetMarshaler(marshaler func(v interface{}) ([]byte, error)) {
	r.marshaler = marshaler
}

// SetUnmarshaler allow to change the unmarshaling function used to decode requests response
func (r *Rehapt) SetUnmarshaler(unmarshaler func(data []byte, v interface{}) error) {
	r.unmarshaler = unmarshaler
}

// SetFail allow to change the function called when TestAssert() encounter an error
func (r *Rehapt) SetFail(fail func(err error)) {
	r.fail = fail
}

// GetVariable allow to retrive a variable value from its name
// nil is returned if variable is not found
func (r *Rehapt) GetVariable(name string) interface{} {
	return r.variables[name]
}

// GetVariable allow to retrive a variable value as a string from its name
// empty string is returned if variable is not found
func (r *Rehapt) GetVariableString(name string) string {
	if value, ok := r.variables[name].(string); ok == true {
		return value
	}
	return ""
}

// SetVariable allow to define manually a variable
func (r *Rehapt) SetVariable(name string, value interface{}) {
	r.variables[name] = value
}

// GetDefaultHeader returns the default request header value from its name
func (r *Rehapt) GetDefaultHeader(name string) string {
	return r.defaultHeaders[name]
}

// SetDefaultHeader allow to add a default request header.
// This header will be added to all requests, however each
// TestCase can override its value by simply adding it to Request Headers
func (r *Rehapt) SetDefaultHeader(name string, value string) {
	r.defaultHeaders[name] = value
}

// SetDefaultTimeDeltaFormat allow to change the default time format
// It is used by TimeDelta, to parse the actual string value as a time
// Default is set to `time.RFC3339` which is ok for JSON
func (r *Rehapt) SetDefaultTimeDeltaFormat(format string) {
	r.defaultTimeDeltaFormat = format
}

// This is the main function of the library
// it executes a given TestCase, i.e. do the request and
// check if the response is matching the expected response
func (r *Rehapt) Test(testcase TestCase) error {
	// If we dont have the minimum, cannot go further.
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
	// Add the testcase defined headers. This override any default header previously set
	for k, v := range testcase.Request.Headers {
		request.Header.Set(k, v)
	}

	// Now execute request
	recorder := httptest.NewRecorder()
	r.httpHandler.ServeHTTP(recorder, request)

	// And start to check result.
	// First check HTTP response code
	if testcase.Response.Code != recorder.Code {
		return fmt.Errorf("response code does not match. Expected %d, got %d", testcase.Response.Code, recorder.Code)
	}

	// Check headers, but not all of them. Only the one expected by the user
	for header, value := range testcase.Response.Headers {
		if value != recorder.Header().Get(header) {
			return fmt.Errorf("response header %v does not match. Expected %v, got %v", header, value, recorder.Header().Get(header))
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

// TestAssert works exactly like Test except it report the error if not nil
// using the Fail function defined by SetFail()
func (r *Rehapt) TestAssert(testcase TestCase) {
	if err := r.Test(testcase); err != nil {
		if r.fail != nil {
			r.fail(err)
		}
	}
}

var variableStoreRegexp = regexp.MustCompile(`^\$([a-zA-Z0-9]+)\$$`)
var variableLoadRegexp = regexp.MustCompile(`_[a-zA-Z0-9]+_`)

func (r *Rehapt) replaceVars(str string) string {
	return variableLoadRegexp.ReplaceAllStringFunc(str, func(name string) string {
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

func (r *Rehapt) storeIfVariable(expected string, actual string) bool {
	elements := variableStoreRegexp.FindStringSubmatch(expected)
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
				r.storeIfVariable("$"+varname+"$", elements[groupid])
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
				return fmt.Errorf("different map sizes. Expected %v, got %v. Expected %v got %v", expectedValue.Len(), actualValue.Len(), expected, actual)
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
			// Compare actual with the loaded value
			r.variables[expectedStr] = actual
			return nil
		}

		if expectedType == reflect.TypeOf(LoadVar("")) {
			// Compare actual with the loaded value which might not be string
			value := r.variables[expectedStr]
			return r.compare(value, actual)
		}

		if actualKind != reflect.String {
			return fmt.Errorf("different kinds. Expected %v, got %v", expectedKind, actualKind)
		}

		actualStr := actualValue.String()

		if r.storeIfVariable(expectedStr, actualStr) == true {
			// This was a variable store operation. no comparison to do
			return nil
		}

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