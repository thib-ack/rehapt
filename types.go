package rehapt

import (
	"io"
	"reflect"
	"time"
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
	Method                       string
	Path                         string
	Headers                      H
	Body                         interface{}
	RawBody                      io.Reader
	NoPathVariableReplacement    bool
	NoHeadersVariableReplacement bool
	NoRawBodyVariableReplacement bool
}

// TestResponse describe the response expected
type TestResponse struct {
	Headers interface{}
	Code    int
	Object  interface{}
	RawBody interface{}
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
// It allow to expect a list of element but without the constraint of order matching
type UnsortedS []interface{}

// StoreVar allow to store the actual value in a variable instead of checking its content
type StoreVar string

// LoadVar allow to load the value of the variable and then compare with actual value
type LoadVar string

// Regexp allow to do advanced regexp expectation.
// If the regexp is invalid, an error is reported.
// If the actual value to compare with is not a string, an error is reported.
// If the actual value does not match the regexp, an error is reported
type Regexp string

// RegexpVars is a mix between Regexp and StoreVar.
// It check if the actual value matches the regexp.
// but all the groups defined in the regexp can be extracted to variables for later reuse
// The Vars hold the mapping groupid: varname.
// For example with Regexp: `^Hello (.*) !$` and Vars: map[int]string{0: "all", 1: "name"}
// then if the actual value is "Hello john !", it will match and 2 vars will be stored:
//  "all" = "Hello john !"  (group 0 is the full match)
//  "name" = "John"
type RegexpVars struct {
	Regexp string
	Vars   map[int]string
}

// Not means we don't expect the given value
// it works as a boolean 'not' operator on the comparison
type Not struct {
	Value interface{}
}

// NumberDelta allow to expect number value with a given +/- delta.
// Delta is compared to math.Abs(expected - actual) which explain why
// if your expected value is 10 with a delta of 3, actual value will match from 7 to 13.
type NumberDelta struct {
	Value float64
	Delta float64
}

// TimeDelta allow to expect time value with a given +/- delta.
// Delta is compared to math.Abs(expected - actual) which explain why
// if your expected time is time.Now() with a delta of 10sec,
// actual value will match from now-10sec to now+10sec
type TimeDelta struct {
	Time   time.Time
	Delta  time.Duration
	Format string
}

type any string

// Any allow you to ignore completely the field
const Any = any("{Any}")

// AnyCode allow you to ignore completely the response code
const AnyCode = -1

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
