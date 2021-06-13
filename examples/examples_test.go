package examples

import (
	. "github.com/thib-ack/rehapt"
	"net/http"
	"testing"
	"time"
)

func TestExampleSimple(t *testing.T) {
	r := setupRehapt(t)

	// Each testcase consist of a description of the request to execute
	// and a description of the expected response
	// By default the response description is exhaustive. if a actual response field is not listed, an error will be triggered.
	// of course if an expected field described here is not present in response, an error will be triggered too.
	r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/user",
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"id":   "55",
				"name": "John",
				"age":  51,
				"pets": S{ // S for slice, M for map. Easy right ?
					M{
						"id":   "123",
						"name": "Pepper the cat",
						"type": "cat",
					},
				},
				"weddingdate": "2019-06-22T16:00:10.123Z",
			},
		},
	})
}

func TestExampleChained(t *testing.T) {
	r := setupRehapt(t)

	// Now imagine you need to call a first API to retrieve some ID in the response body
	// and then use this ID to call a second API
	r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/user",
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"id":   "55",
				"name": "John",
				"age":  51,
				"pets": S{ // S for slice, M for map. Easy right ?
					M{
						"id":   "$catid$",        // Note ! Here we don't compare but register the value returned here inside a variable named catid,
						"name": "Pepper the cat", // this form is a shortcut. You could also use `StoreVar("catid")`
						"type": "cat",
					},
				},
				"weddingdate": "2019-06-22T16:00:10.123Z",
			},
		},
	})

	r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/cat/_catid_", // Note ! Here we use the previously registered catid variable. $xyz$ for register _xyz_ for use.
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"id":   "123",
				"name": "Pepper the cat",
				"age":  3,
				"owner": M{
					"id":   "55",
					"name": "John",
					"age":  51.0,
				},
				"toys": S{"ball", "plastic mouse"},
			},
		},
	})
}

func TestExampleAdvanced(t *testing.T) {
	r := setupRehapt(t)

	// Now lets demonstrate more cool features
	r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/user",
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"id":   StoreVar("id"),
				"name": "John",
				"age":  StoreVar("age"), // StoreVar works on any actual type (int here) not only strings
				"pets": S{
					PartialM{ // A partial map ignores unlisted field instead of reporting errors
						"id":   "$catid$", // We dont compare but register the cat ID returned here
						"type": Any,       // We ignore this field
					},
				},
				"weddingdate": TimeDelta{Time: time.Date(2019, time.June, 22, 16, 0, 0, 0, time.UTC), Delta: 20 * time.Second}, // We can compare dates with delta
			},
		},
	})

	// You can work with the stored variables if you have additional checks to do.
	if r.GetVariableString("catid") == "" {
		t.Error("missing cat ID")
	}

	// You can also define your own variables
	_ = r.SetVariable("myvar", M{
		"id":   "55",
		"name": "John",
		"age":  51,
	})

	// Or use Test() which return the error if you want to handle it in your test
	err := r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/cat/_catid_", // Here we use the previously registered catid variable. $xyz$ for register _xyz_ for use.
		},
		Response: TestResponse{
			Code:    http.StatusOK,
			Headers: M{"X-Pet-Type": S{"Cat"}}, // Check for header presence in response
			Object: M{
				"id":   "_catid_", // Here we load the previously registred var. If does not match with returned value -> error (try to change in example server)
				"age":  Any,
				"name": RegexpVars{Regexp: `(.*) the cat`, Vars: map[int]string{1: "catname"}}, // We can store vars using regexp groups
				"toys": S{
					"ball",
					Regexp(`^plastic .*$`), // We can compare with regexp
				},
				"owner": LoadVar("myvar"), // Here we load the previously registered var. This works even for maps, slices, etc.
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	if r.GetVariableString("catname") != "Pepper" {
		t.Error("incorrect cat name")
	}
}

func TestExampleInvalidAuth(t *testing.T) {
	r := setupRehapt(t)

	// This override the default Authorization header (set in setupRehapt)
	r.TestAssert(TestCase{
		Request: TestRequest{
			Method:  "GET",
			Path:    "/api/user",
			Headers: H{"Authorization": {"invalid"}},
		},
		Response: TestResponse{
			Code: http.StatusUnauthorized,
			Object: M{
				"error": "unauthorized",
			},
		},
	})
}
