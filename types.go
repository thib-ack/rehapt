package rehapt

import (
	"time"
)

type TestCase struct {
	Request  TestRequest
	Response TestResponse
}

type TestRequest struct {
	Method  string
	Path    string
	Headers H
	Body    interface{}
}

type TestResponse struct {
	Headers H
	Code    int
	Object  interface{}
	Raw     interface{}
}

// Type H is used to quickly define Headers
type H map[string]string

// Type M is used to quickly build a map
type M map[string]interface{}

// Partial Map allow to expect some fields but ignore the un-listed ones instead of reporting missing
type PartialM map[string]interface{}

// Type S is used to quickly build a slice
type S []interface{}

// Unsorted Slice allow to expect a list of element but without the constraint of order matching
type UnsortedS []interface{}

type StoreVar string

type LoadVar string

type Regexp string

type RegexpVars struct {
	Regexp string
	Vars   map[int]string
}

type NumberDelta struct {
	Value float64
	Delta float64
}

type TimeDelta struct {
	Time   time.Time
	Delta  time.Duration
	Format string
}

type any string

const Any = any("{Any}")
