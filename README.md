
# rehapt <a href="https://travis-ci.org/thib-ack/rehapt"><img src="https://travis-ci.org/thib-ack/rehapt.svg?branch=master"></a> <a href="https://goreportcard.com/report/thib-ack/rehapt"><img src="http://goreportcard.com/badge/thib-ack/rehapt"></a> <a href="https://godoc.org/github.com/thib-ack/rehapt"><img src="https://godoc.org/github.com/thib-ack/rehapt?status.svg" alt="GoDoc reference"></a>

Rehapt is a Golang declarative REST HTTP API testing library.
You describe how you expect your HTTP API to behave and the library take care of comparing the expected and actual response elements.

This library has been designed to work very well for JSON APIs but it can 
handle any other format as long as it marshal/unmarshal to/from simple maps and slices

## Features

- Support Go v1.7+
- Work on `http.Handler`, no need to really start your http server
- Describe the request you want to perform
- Describe the object you expect as response
- Easy to configure
- Define default headers (useful for auth for example)
- Include a variable system to extract values from API responses and reuse them later
- Smart response expectations:
    - Ignore some fields
    - Regexp string expectation
    - Expect times with delta
    - Expect numbers with delta
    - Store fields in variables
    - Load variables previously stored
    - Unsorted slice expectation if order doesn't matter
    - Partial map expectation when only some keys matter
- No third-party dependencies

## Installation

```bash
go get github.com/thib-ack/rehapt
```

## Examples

See [examples](https://github.com/thib-ack/rehapt/blob/master/examples) folder for more examples.

#### Simple example

```go
package example

import (
    . "github.com/thib-ack/rehapt"
    "net/http"
    "testing"
)

func TestAPISimple(t *testing.T) {
    r := NewRehapt(yourHttpServerMux)

    // Each testcase consist of a description of the request to execute
    // and a description of the expected response
    // By default the response description is exhaustive. 
    // If an actual response field is not listed here, an error will be triggered
    // of course if an expected field described here is not present in response, an error will be triggered too.
    r.TestAssert(TestCase{
        Request: TestRequest{
            Method: "GET",
            Path:   "/api/user/1",
        },
        Response: TestResponse{
            Code: http.StatusOK,
            Object: M{
                "id":   "1",
                "name": "John",
                "age":  51,
                "pets": S{ // S for slice, M for map. Easy right ?
                    M{
                        "id":   "2",
                        "name": "Pepper the cat",
                        "type": "cat",
                    },
                },
                "weddingdate": "2019-06-22T16:00:00.000Z",
            },
        },
    })
}
```

#### Advanced examples

Obviously more advanced features are supported:

```go
package example

import (
    . "github.com/thib-ack/rehapt"
    "net/http"
    "testing"
    "time"
)

func TestAPIAdvanced(t *testing.T) {
    r := NewRehapt(yourHttpServerMux)

    r.TestAssert(TestCase{
        Request: TestRequest{
            Method: "GET",
            // Add headers to request. (default headers are also supported)
            Headers: H{"X-Custom": "my value"}, 
            Path:   "/api/user/1",
        },
        Response: TestResponse{
            Code: http.StatusOK,
            // Check for headers presence in response
            Headers: H{"X-Pet-Type": "Cat"},
            Object: M{
                "id":   "1",
                // We can ignore a specific field
                "name": Any,
                // We can expect numbers with a given delta
                "age":  NumberDelta{Value: 50, Delta: 10},
                "pets": S{
                    // We can expect a partial map. 
                    // the keys not listed here will be ignored instead of returned as missing
                    PartialM{
                        "id":   "2",
                        // We can expect with regexp
                        "name": Regexp(`[A-Za-z]+ the cat`),
                        "type": "cat",
                        // We can expect a slice without order constraint
                        "toys": UnsortedS{"mouse", "ball"},
                    },
                },
                // We can expect dates with a given delta
                "weddingdate": TimeDelta{Time: time.Now(), Delta: 24 * time.Hour},
            },
        },
    })
}
```

Rehapt also includes a variable system, used to extract values from API responses and use them in later API calls.

```go
package example

import (
    "fmt"
    . "github.com/thib-ack/rehapt"
    "net/http"
    "testing"
)

func TestAPIVariables(t *testing.T) {
    r := NewRehapt(yourHttpServerMux)

    r.TestAssert(TestCase{
        Request: TestRequest{
            Method: "GET",
            Path:   "/api/user/1",
        },
        Response: TestResponse{
            Code: http.StatusOK,
            Object: M{
                // StoreVar doesn't check the actual value but store it in a variable named "age" here
                "age":  StoreVar("age"),
                "pets": S{
                    M{
                        // This shortcut act like a StoreVar("catid")
                        "id":   "$catid$",
                    },
                },
            },
        },
    })

    // And we can reuse the variables in a next API call
    r.TestAssert(TestCase{
        Request: TestRequest{
            Method: "GET",
            // Here we indicate to use the variable catid value. 
            // for example this will call /api/cat/2 if value in previous request was 2
            Path:   "/api/cat/_catid_",
        },
        Response: TestResponse{
            Code: http.StatusOK,
            Object: M{
                // LoadVar load the variable value and check with actual response value.
                // Here it will report an error if cat's age is not the same as John's age
                // which were extracted from previous request
                "age":  LoadVar("age"),
            },
        },
    })

    // Or we can play with the variables if we need
    fmt.Println("Its age is ", r.GetVariable("age"))
    // We can also define them by hand. Any type of value is supported
    r.SetVariable("myvar", M{"key": "value"})
}
```

## License

MIT - Thibaut Ackermann
