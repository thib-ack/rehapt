package rehapt

import (
	"fmt"
	"reflect"
)

// ErrorHandler is the interface used to report errors when found by TestAssert().
// Note that *testing.T implements this interface
type ErrorHandler interface {
	Errorf(format string, args ...interface{})
}

// TestCase is the base type supported to describe a test.
// It is the object taken as parameters in Test() and TestAssert()
type TestCase struct {
	Request  TestRequest
	Response TestResponse
}

// TestRequest describe the request to be executed
type TestRequest struct {
	Method        string
	Path          interface{}
	Headers       H
	Body          interface{}
	BodyMarshaler MarshalFn
}

// TestResponse describe the response expected
type TestResponse struct {
	Headers         interface{}
	Code            interface{}
	Body            interface{}
	BodyUnmarshaler UnmarshalFn
}

// H declare a Headers map.
// It is used to quickly define Headers within your requests
type H map[string][]string

// M declare a Map.
// It is used to quickly build a map within your expected response body
type M map[string]interface{}

// PartialM declare a Partial Map.
// It is used to expect some fields but ignore the un-listed ones instead of reporting missing
type PartialM map[string]interface{}

// S declare a Slice.
// It is used to quickly build a slice within your expected response body
type S []interface{}

// UnsortedS declare an Unsorted Slice.
// It allows to expect a list of element but without the constraint of order matching
type UnsortedS []interface{}

type CompareFn func(r *Rehapt, ctx compareCtx) error

type ReplaceFn func(r *Rehapt) (string, error)

type MarshalFn func(v interface{}) ([]byte, error)

func RawMarshaler(v interface{}) ([]byte, error) {
	if s, ok := v.(string); ok == true {
		return []byte(s), nil
	} else if b, ok := v.([]byte); ok == true {
		return b, nil
	} else {
		return nil, fmt.Errorf("only string or []byte supported")
	}
}

type UnmarshalFn func(data []byte, v interface{}) error

func RawUnmarshaler(data []byte, out interface{}) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("out should be a non-nil pointer")
	}

	if rv.IsValid() {
		pv := rv.Elem()
		pv.Set(reflect.ValueOf(string(data)))
	}
	return nil
}

type compareCtx struct {
	Expected      interface{}
	ExpectedKind  reflect.Kind
	ExpectedType  reflect.Type
	ExpectedValue reflect.Value
	Actual        interface{}
	ActualKind    reflect.Kind
	ActualType    reflect.Type
	ActualValue   reflect.Value
}

type comparator struct {
	ExpectedKind reflect.Kind
	ExpectedType reflect.Type
	Compare      func(compareCtx) error
}
