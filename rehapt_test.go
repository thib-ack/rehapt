package rehapt_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "rehapt"
	"testing"
	"time"
)

type testContext struct {
	r      *Rehapt
	server *http.ServeMux
}

func setupTest(t *testing.T) *testContext {
	server := http.NewServeMux()

	c := &testContext{
		r:      NewRehapt(server),
		server: server,
	}
	c.r.SetFail(func(err error) {
		t.Error(err)
	})
	return c
}

// small helper to expect an error or report failure is no error or incorrect error
func ExpectError(err error, str string) string {
	if err == nil {
		return fmt.Sprintf("Expected '%v', got no error", str)
	}
	e := err.Error()
	if e == str {
		// That's OK, error match expected
		return ""
	}
	return fmt.Sprintf("Expected '%v', got '%v'", str, e)
}

// small helper to expect no error or report failure if error
func ExpectNil(err error) string {
	if err == nil {
		// That's OK
		return ""
	}
	e := err.Error()
	return fmt.Sprintf("Expected no error, got '%v'", e)
}

// Now finally our tests
// Begin with valid cases

func TestValidCaseSimpleString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"ok"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: "ok",
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleBool(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `true`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: true,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleInt(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `10`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: 10,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: uint(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: 10.0,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: int64(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleFloat(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `10.0`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: 10.0,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: uint(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: 10,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: int64(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleMap(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name": "John", "Age": 51}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"name": "John",
				"Age":  51,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleSlice(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `["John", "Doe", 99]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: S{
				"John",
				"Doe",
				99,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleRequestBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		var body struct {
			Msg string `json:"msg"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if expected, actual := body.Msg, "ok"; expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body: M{
				"msg": "ok",
			},
		},
		Response: TestResponse{
			Code:   http.StatusAccepted,
			Object: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleRequestDefaultHeader(t *testing.T) {
	c := setupTest(t)

	c.r.SetDefaultHeader("X-Custom", "custom value 123")

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := req.Header.Get("X-Custom"), "custom value 123"; expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if actual, expected := c.r.GetDefaultHeader("X-Custom"), "custom value 123"; actual != expected {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestValidCaseSimpleRequestHeader(t *testing.T) {
	c := setupTest(t)

	c.r.SetDefaultHeader("X-Custom", "default value")

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := req.Header.Get("X-Custom"), "custom value 123"; expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:  "POST",
			Path:    "/api/test",
			Headers: H{"X-Custom": "custom value 123"},
			Body:    nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleResponseHeader(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Custom", "custom value 123")
		w.WriteHeader(http.StatusOK)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Headers: H{
				"X-Custom": "custom value 123",
			},
			Object: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleRawString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  "Hello this is plain text",
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleRawStringVar(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  "$body$",
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := c.r.GetVariable("body"), "Hello this is plain text"; expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestValidCaseSimpleRawRegexp(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  Regexp(`^H[a-z ]+ [0-9]+$`),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseSimpleRawRegexpVars(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw: RegexpVars{
				Regexp: `^H[a-z ]+ ([0-9]+)$`,
				Vars:   map[int]string{1: "counter"},
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := c.r.GetVariable("counter"), "1234"; expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestValidCaseSimpleAssert(t *testing.T) {
	c := setupTest(t)

	c.r.SetFail(func(err error) {
		t.Errorf("this function should not be called")
	})

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"ok"`)
	})

	c.r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: "ok",
		},
	})
}

func TestValidCaseAdvancedIgnore(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": Any,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedStoreVarString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": StoreVar("myvar"),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "150 - high - end", c.r.GetVariableString("myvar"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestValidCaseAdvancedStoreVarNumber(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": 1580}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": StoreVar("myvar"),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := float64(1580), c.r.GetVariable("myvar"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestValidCaseAdvancedLoadVarString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "value"}`)
	})

	c.r.SetVariable("myvar", "value")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": LoadVar("myvar"),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedLoadVarNumber(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": 1580}`)
	})

	c.r.SetVariable("myvar", 1580)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": LoadVar("myvar"),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedRegexp(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": Regexp(`^[0-9]{3} - .* - end$`),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedRegisterVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "high"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": "$stats$",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := c.r.GetVariable("stats"), "high"; expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestValidCaseAdvancedUseVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok"}`)
	})

	c.r.SetVariable("id", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test/_id_",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"status": "ok",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedPartialMap(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"A": 1, "B": 2, "C": 3}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: PartialM{"A": 1},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedUnsortedSlice(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `["A", "B", "C"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: UnsortedS{"B", "C", "A"},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedNumberDeltaExact(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: NumberDelta{
				Value: 555,
				Delta: 0,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedNumberDeltaLower(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: NumberDelta{
				Value: 550,
				Delta: 5,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedNumberDeltaGreater(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: NumberDelta{
				Value: 560,
				Delta: 5,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedTimeDeltaExact(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
				Delta: 0,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedTimeDeltaBefore(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 30, 0, time.UTC),
				Delta: 1 * time.Second,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedTimeDeltaAfter(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 31, 0, time.UTC),
				Delta: 1 * time.Second,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedTimeDeltaCustomDefaultFormat(t *testing.T) {
	c := setupTest(t)

	c.r.SetDefaultTimeDeltaFormat("Day 2006-01-02 Hour 15:04:05Z07:00")

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"Day 2020-04-11 Hour 20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
				Delta: 0,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedTimeDeltaCustomFormat(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"Day 2020-04-11 Hour 20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:   time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
				Delta:  0,
				Format: "Day 2006-01-02 Hour 15:04:05Z07:00",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestValidCaseAdvancedRegexpVars(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"The output value is: Hello and World."`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: RegexpVars{
				Regexp: `.*: (\w+) and (\w+)\.`,
				Vars:   map[int]string{1: "first", 2: "second"},
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := c.r.GetVariable("first"), "Hello"; expected != actual {
		t.Errorf("expected value %v but got %v, ", expected, actual)
	}

	if expected, actual := c.r.GetVariable("second"), "World"; expected != actual {
		t.Errorf("expected value %v but got %v, ", expected, actual)
	}
}

func TestValidCaseAdvancedRegexpVarsOnlyFullMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"--header--content--footer--"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: RegexpVars{
				Regexp: `--header--.+--footer--`,
				Vars:   map[int]string{0: "full"},
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := c.r.GetVariable("full"), "--header--content--footer--"; expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

// And now invalid cases

func TestInvalidCaseNilMarshaler(t *testing.T) {
	c := setupTest(t)

	c.r.SetMarshaler(nil)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `nil marshaler`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseNilUnmarshaler(t *testing.T) {
	c := setupTest(t)

	c.r.SetUnmarshaler(nil)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `nil unmarshaler`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseNilHTTPHandler(t *testing.T) {
	c := setupTest(t)

	c.r.SetHttpHandler(nil)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `nil HTTP handler`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseMissingHTTPMethod(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `incomplete testcase. Missing HTTP method`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseInvalidHTTPMethod(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "NOT CORRECT",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `failed to build HTTP request. net/http: invalid method "NOT CORRECT"`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseMissingURLPath(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `incomplete testcase. Missing URL path`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseCannotMarshalRequestBody(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   M{"n": json.Number(`invalid`)}, // This is refused by json.Marshal
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `failed to marshal the testcase request body. json: invalid number literal "invalid"`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectHTTPCode(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `response code does not match. Expected 200, got 401`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseAssert(t *testing.T) {
	c := setupTest(t)

	called := false
	c.r.SetFail(func(err error) {
		called = true
	})

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"ok"`)
	})

	c.r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: "not ok",
		},
	})

	if called == false {
		t.Errorf("Fail function should have been called")
	}
}

func TestInvalidCaseIncorrectObjectString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"ok"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: "nok",
		},
	})

	if e := ExpectError(err, `strings does not match. Expected 'nok', got 'ok'`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectType(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/string", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"ok"`)
	})
	c.server.HandleFunc("/api/int", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `1`)
	})
	c.server.HandleFunc("/api/float", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `1.0`)
	})
	c.server.HandleFunc("/api/bool", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `true`)
	})
	c.server.HandleFunc("/api/map", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"msg": "ok"}`)
	})
	c.server.HandleFunc("/api/slice", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `["ok"]`)
	})

	type unexportedStruct struct{}

	tests := []struct {
		Path   string
		Object interface{}
		Error  string
	}{
		// Int
		{Path: "string", Object: 1, Error: "different kinds. Expected int, got string"},
		{Path: "bool", Object: 1, Error: "different kinds. Expected int, got bool"},
		{Path: "map", Object: 1, Error: "different kinds. Expected int, got map"},
		{Path: "slice", Object: 1, Error: "different kinds. Expected int, got slice"},
		// Uint
		{Path: "string", Object: uint(1), Error: "different kinds. Expected uint, got string"},
		{Path: "bool", Object: uint(1), Error: "different kinds. Expected uint, got bool"},
		{Path: "map", Object: uint(1), Error: "different kinds. Expected uint, got map"},
		{Path: "slice", Object: uint(1), Error: "different kinds. Expected uint, got slice"},
		// String
		{Path: "int", Object: "ok", Error: "different kinds. Expected string, got float64"}, // TODO can we force json.Unmarshal to use int ?
		{Path: "float", Object: "ok", Error: "different kinds. Expected string, got float64"},
		{Path: "bool", Object: "ok", Error: "different kinds. Expected string, got bool"},
		{Path: "map", Object: "ok", Error: "different kinds. Expected string, got map"},
		{Path: "slice", Object: "ok", Error: "different kinds. Expected string, got slice"},
		// Float32
		{Path: "string", Object: float32(0.1), Error: "different kinds. Expected float32, got string"},
		{Path: "bool", Object: float32(0.1), Error: "different kinds. Expected float32, got bool"},
		{Path: "map", Object: float32(0.1), Error: "different kinds. Expected float32, got map"},
		{Path: "slice", Object: float32(0.1), Error: "different kinds. Expected float32, got slice"},
		// Float64
		{Path: "string", Object: float64(0.1), Error: "different kinds. Expected float64, got string"},
		{Path: "bool", Object: float64(0.1), Error: "different kinds. Expected float64, got bool"},
		{Path: "map", Object: float64(0.1), Error: "different kinds. Expected float64, got map"},
		{Path: "slice", Object: float64(0.1), Error: "different kinds. Expected float64, got slice"},
		// Bool
		{Path: "string", Object: true, Error: "different kinds. Expected bool, got string"},
		{Path: "int", Object: true, Error: "different kinds. Expected bool, got float64"},
		{Path: "float", Object: true, Error: "different kinds. Expected bool, got float64"},
		{Path: "map", Object: true, Error: "different kinds. Expected bool, got map"},
		{Path: "slice", Object: true, Error: "different kinds. Expected bool, got slice"},
		// Map
		{Path: "string", Object: M{}, Error: "different kinds. Expected map, got string"},
		{Path: "int", Object: M{}, Error: "different kinds. Expected map, got float64"},
		{Path: "float", Object: M{}, Error: "different kinds. Expected map, got float64"},
		{Path: "bool", Object: M{}, Error: "different kinds. Expected map, got bool"},
		{Path: "slice", Object: M{}, Error: "different kinds. Expected map, got slice"},
		// Slice
		{Path: "string", Object: S{}, Error: "different kinds. Expected slice, got string"},
		{Path: "int", Object: S{}, Error: "different kinds. Expected slice, got float64"},
		{Path: "float", Object: S{}, Error: "different kinds. Expected slice, got float64"},
		{Path: "bool", Object: S{}, Error: "different kinds. Expected slice, got bool"},
		{Path: "map", Object: S{}, Error: "different kinds. Expected slice, got map"},
		// Struct
		{Path: "string", Object: struct{}{}, Error: "unexpected struct type struct {}"},
		{Path: "int", Object: struct{}{}, Error: "unexpected struct type struct {}"},
		{Path: "float", Object: struct{}{}, Error: "unexpected struct type struct {}"},
		{Path: "bool", Object: struct{}{}, Error: "unexpected struct type struct {}"},
		{Path: "slice", Object: struct{}{}, Error: "unexpected struct type struct {}"},
		// Unhandled
		{Path: "string", Object: complex(1, 2), Error: "unhandled type complex128"},
	}

	for _, test := range tests {
		err := c.r.Test(TestCase{
			Request: TestRequest{
				Method: "GET",
				Path:   "/api/" + test.Path,
				Body:   nil,
			},
			Response: TestResponse{
				Code:   http.StatusOK,
				Object: test.Object,
			},
		})

		if e := ExpectError(err, test.Error); e != "" {
			t.Error(e)
		}
	}
}

func TestInvalidCaseIncorrectResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		// This is not valid JSON
		fmt.Fprintf(w, `{"error": invalid...`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: Any,
		},
	})

	if e := ExpectError(err, `cannot unmarshal response body. invalid character 'i' looking for beginning of value`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseHeader(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Custom", "not right value")
		w.WriteHeader(http.StatusOK)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Headers: H{
				"X-Custom": "custom value 123",
			},
			Object: nil,
		},
	})

	if e := ExpectError(err, `response header X-Custom does not match. Expected custom value 123, got not right value`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseNilResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: "anything",
		},
	})

	if e := ExpectError(err, `expected anything but got nil`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseResponseBodyExpectedNil(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"success"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: nil,
		},
	})

	if e := ExpectError(err, `expected is nil but got success`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseSliceDifferentSize(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `["A", "B"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: S{"A"},
		},
	})

	if e := ExpectError(err, `different slice sizes. Expected 1, got 2. Expected [A] got [A B]`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseSliceElementDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `["A", "B"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: S{"A", "C"},
		},
	})

	if e := ExpectError(err, `slice element 1 does not match. strings does not match. Expected 'C', got 'B'`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseMapDifferentKeyType(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: map[int]string{1: "test"},
		},
	})

	if e := ExpectError(err, `different map key types. Expected int, got string`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseMapDifferentSize(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: M{"key": "value", "foo": "bar"},
		},
	})

	if e := ExpectError(err, `different map sizes. Expected 2, got 1. Expected map[foo:bar key:value] got map[key:value]`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseMapKeyNotFound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: M{"foo": "bar"},
		},
	})

	if e := ExpectError(err, `expected key foo not found`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectResponseMapElementDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: M{"key": "bar"},
		},
	})

	if e := ExpectError(err, `map element [key] does not match. strings does not match. Expected 'bar', got 'value'`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectNumberDeltaNotNumber(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"hi"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: NumberDelta{
				Value: 500,
				Delta: 10,
			},
		},
	})

	if e := ExpectError(err, `different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectNumberDeltaLower(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `500`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: NumberDelta{
				Value: 450,
				Delta: 49,
			},
		},
	})

	if e := ExpectError(err, `max difference between 450 and 500 allowed is 49, but difference was 50`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectNumberDeltaGreater(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `500`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: NumberDelta{
				Value: 550,
				Delta: 49,
			},
		},
	})

	if e := ExpectError(err, `max difference between 550 and 500 allowed is 49, but difference was 50`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectTimeDeltaNotString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `1000`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 31, 0, time.UTC),
				Delta: 1 * time.Second,
			},
		},
	})

	if e := ExpectError(err, `different kinds. Expected string, got float64`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectTimeDeltaNotTime(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"hello"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 31, 0, time.UTC),
				Delta: 1 * time.Second,
			},
		},
	})

	if e := ExpectError(err, `invalid time. parsing time "hello" as "2006-01-02T15:04:05Z07:00": cannot parse "hello" as "2006"`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectTimeDeltaBefore(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 29, 0, time.UTC),
				Delta: 1 * time.Second,
			},
		},
	})

	if e := ExpectError(err, `max difference between 2020-04-11 20:10:29 +0000 UTC and 2020-04-11 20:10:30.123 +0000 UTC allowed is 1s, but difference was -1.123s`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseIncorrectTimeDeltaAfter(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: TimeDelta{
				Time:  time.Date(2020, time.April, 11, 20, 10, 32, 0, time.UTC),
				Delta: 1 * time.Second,
			},
		},
	})

	if e := ExpectError(err, `max difference between 2020-04-11 20:10:32 +0000 UTC and 2020-04-11 20:10:30.123 +0000 UTC allowed is 1s, but difference was 1.877s`); e != "" {
		t.Error(e)
	}
}

func TestValidCaseUnsortedSliceElementNotFound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `["A", "B", "C"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: UnsortedS{"B", "C", "E"},
		},
	})

	if e := ExpectError(err, `expected element E at index 2 not found`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRegexpFailParsing(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": Regexp(`^[0-9](3 - .* - end$`),
			},
		},
	})

	if e := ExpectError(err, "map element [stats] does not match. error parsing regexp: missing closing ): `^[0-9](3 - .* - end$`"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRegexpDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": Regexp(`^[a-z]{3} - .* - end$`),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. regexp '^[a-z]{3} - .* - end$' does not match '150 - high - end'`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseBoolDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `true`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:   http.StatusOK,
			Object: false,
		},
	})

	if e := ExpectError(err, `bools does not match. Expected false, got true`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRegexpVarsNotString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `1000`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: RegexpVars{
				Regexp: `^([0-9]{3})$`,
				Vars:   nil,
			},
		},
	})

	if e := ExpectError(err, `different kinds. Expected string, got float64`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRegexpVarsFailParsing(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": RegexpVars{
					Regexp: `^[0-9](3 - .* - end$`,
					Vars:   nil,
				},
			},
		},
	})

	if e := ExpectError(err, "map element [stats] does not match. error parsing regexp: missing closing ): `^[0-9](3 - .* - end$`"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRegexpVarsDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": RegexpVars{
					Regexp: `^[a-z]{3} - (.*) - end$`,
					Vars:   map[int]string{1: "v1"},
				},
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. regexp '^[a-z]{3} - (.*) - end$' does not match '150 - high - end'`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRegexpVarsOverflowIndexIgnored(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Object: M{
				"stats": RegexpVars{
					Regexp: `^[0-9]{3} - (.*) - end$`,
					Vars:   map[int]string{2: "v1"},
				},
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. expected variable index 2 overflow regexp group count of 2`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawUnhandled(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  1234,
		},
	})

	if e := ExpectError(err, "unsupported Raw object type int"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawStringDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  "Hello this is plain text",
		},
	})

	if e := ExpectError(err, "response body does not match. Expected Hello this is plain text, got Hello this is plain text 1234"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawRegexpFailParsing(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  Regexp(`^H[a-z ]+ ([0-9]+$`),
		},
	})

	if e := ExpectError(err, "error parsing regexp: missing closing ): `^H[a-z ]+ ([0-9]+$`"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawRegexpDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw:  Regexp(`^H[a-z ]+ [0-9]$`),
		},
	})

	if e := ExpectError(err, "regexp '^H[a-z ]+ [0-9]$' does not match 'Hello this is plain text 1234'"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawRegexpVarsFailParsing(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw: RegexpVars{
				Regexp: `^H[a-z ]+ ([0-9]+$`,
				Vars:   map[int]string{1: "counter"},
			},
		},
	})

	if e := ExpectError(err, "error parsing regexp: missing closing ): `^H[a-z ]+ ([0-9]+$`"); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawRegexpVarsDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw: RegexpVars{
				Regexp: `^H[a-z ]+ ([0-9])$`,
				Vars:   map[int]string{1: "counter"},
			},
		},
	})

	if e := ExpectError(err, `regexp '^H[a-z ]+ ([0-9])$' does not match 'Hello this is plain text 1234'`); e != "" {
		t.Error(e)
	}
}

func TestInvalidCaseRawRegexpVarsOverflowIndex(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Raw: RegexpVars{
				Regexp: `^H[a-z ]+ ([0-9]+)$`,
				Vars:   map[int]string{2: "counter"},
			},
		},
	})

	if e := ExpectError(err, `expected variable index 2 overflow regexp group count of 2`); e != "" {
		t.Error(e)
	}
}
