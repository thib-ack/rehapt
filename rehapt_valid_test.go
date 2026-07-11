package rehapt_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	. "github.com/thib-ack/rehapt"
)

type testContext struct {
	r      *Rehapt
	server *http.ServeMux
}

func setupTest(t *testing.T) *testContext {
	server := http.NewServeMux()

	c := &testContext{
		r:      NewRehapt(t, server),
		server: server,
	}

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

// small helper to make sure the Errorf function is called
type testingT struct {
	called bool
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.called = true
}

// Now finally our tests

// All valid test cases

func TestOKStringResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"ok"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: "ok",
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKBoolResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `true`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: true,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKIntResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `10`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: 10,
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
			Code: http.StatusOK,
			Body: uint(10),
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
			Code: http.StatusOK,
			Body: 10.0,
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
			Code: http.StatusOK,
			Body: int64(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKFloatResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `10.0`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: 10.0,
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
			Code: http.StatusOK,
			Body: uint(10),
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
			Code: http.StatusOK,
			Body: 10,
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
			Code: http.StatusOK,
			Body: int64(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNotResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `10`)
	})
	c.server.HandleFunc("/api/test-str", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"hello"`)
	})
	c.server.HandleFunc("/api/test-empty", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ``)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Not(15),
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
			Code: http.StatusOK,
			Body: Not("world"),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-str",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Not("world"),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-str",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Not(false),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-str",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Not(nil),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-str",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Not(10.0),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-empty",
			Body:   nil,
		},
		Response: TestResponse{
			Code: Any(),
			Body: Not("hello"),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-empty",
			Body:   nil,
		},
		Response: TestResponse{
			Code: Any(),
			Body: Not(false),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test-empty",
			Body:   nil,
		},
		Response: TestResponse{
			Code: Any(),
			Body: Not(10),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKAndResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"hello"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: And("hello", Regexp("h...o")),
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
			Code: http.StatusOK,
			Body: And(StoreVar("v"), "hello"),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if c.r.GetVariableString("v") == "" {
		t.Error("missing cat ID")
	}
}

func TestOKOrResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"hello"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Or(), // empty or
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
			Code: http.StatusOK,
			Body: Or("hello", "world"),
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
			Code: http.StatusOK,
			Body: Or("world", "hello"),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKMapResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"name": "John", "Age": 51}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"name": "John",
				"Age":  51,
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKPartialMapResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"name": "John", "Age": 51}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: PartialM{
				"name": "John",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKSliceResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["John", "Doe", 99]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: S{
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

func TestOKUnsortedSliceResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["John", "Doe", 99]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: UnsortedS{
				"Doe",
				99,
				"John",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKTimeDateResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"date": "2020-04-11T20:10:30.123Z"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"date": time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestPathLoadVarShortcut(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	_ = c.r.SetVariable("catid", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/_catid_",
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestPathNoReplacement(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/_catid_", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	_ = c.r.SetVariable("catid", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   NoReplacement("/api/_catid_"),
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestPathInvalidType(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   1234,
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectError(err, "invalid path type int, only string or rehapt.ReplaceFn supported"); e != "" {
		t.Error(e)
	}
}

func TestOKRequestBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		var body struct {
			Msg string `json:"msg"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if expected, actual := "ok", body.Msg; expected != actual {
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
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestRawBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if expected, actual := "This is a raw plain/text body", string(body); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:        "POST",
			Path:          "/api/test",
			BodyMarshaler: RawMarshaler,
			Body:          "This is a raw plain/text body",
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method:        "POST",
			Path:          "/api/test",
			BodyMarshaler: RawMarshaler,
			Body:          []byte("This is a raw plain/text body"),
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestRawBodyLoadVarShortcut(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if expected, actual := "The cat 123 won", string(body); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
	})

	_ = c.r.SetVariable("catid", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:        "POST",
			Path:          "/api/test",
			BodyMarshaler: RawMarshaler,
			Body:          c.r.ReplaceVars("The cat _catid_ won"),
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestRawLoadVarShortcutBodyNoReplacement(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if expected, actual := "The cat _catid_ won", string(body); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
	})

	_ = c.r.SetVariable("catid", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:        "POST",
			Path:          "/api/test",
			BodyMarshaler: RawMarshaler,
			Body:          "The cat _catid_ won",
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKDefaultRequestHeader(t *testing.T) {
	c := setupTest(t)

	// Set the default header (will be set for all requests)
	c.r.SetDefaultHeader("X-Custom", "custom value 123")

	// We can check its value too
	if actual, expected := c.r.GetDefaultHeader("X-Custom"), "custom value 123"; actual != expected {
		t.Errorf("expected value %v but got %v", expected, actual)
	}

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := "custom value 123", req.Header.Get("X-Custom"); expected != actual {
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
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKDefaultRequestHeaders(t *testing.T) {
	c := setupTest(t)

	// Set the default headers (will be set for all requests)
	c.r.SetDefaultHeaders(http.Header{"X-Custom": []string{"custom value 123"}})

	// We can check its value too
	headers := c.r.GetDefaultHeaders()
	if actual, expected := headers.Get("X-Custom"), "custom value 123"; actual != expected {
		t.Errorf("expected value %v but got %v", expected, actual)
	}

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := "custom value 123", req.Header.Get("X-Custom"); expected != actual {
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
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestHeader(t *testing.T) {
	c := setupTest(t)

	// We set a default header, but it will be overloaded by the request one
	c.r.SetDefaultHeader("X-Custom", "default value")

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := "custom value 123", req.Header.Get("X-Custom"); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:  "POST",
			Path:    "/api/test",
			Headers: H{"X-Custom": {"custom value 123"}},
			Body:    nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	// Now we don't overload the header, the default one should be used again
	// This test is important because we had an issue where the default header was erased when overloaded
	c.server.HandleFunc("/api/test2", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := "default value", req.Header.Get("X-Custom"); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
		w.WriteHeader(http.StatusOK)
	})

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test2",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestHeaderGetVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := "123", req.Header.Get("X-Header"); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
		w.WriteHeader(http.StatusOK)
	})

	_ = c.r.SetVariable("hdr", "X-Header")
	_ = c.r.SetVariable("catid", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Headers: H{
				c.r.GetVariableString("hdr"): {c.r.GetVariableString("catid")},
			},
		},
		Response: TestResponse{
			Code: http.StatusOK,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRequestHeaderNoReplacement(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		if expected, actual := "_catid_", req.Header.Get("_hdr_"); expected != actual {
			t.Errorf("expected value %v but got %v", expected, actual)
		}
		w.WriteHeader(http.StatusOK)
	})

	_ = c.r.SetVariable("hdr", "X-Header")
	_ = c.r.SetVariable("catid", "123")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:  "POST",
			Path:    "/api/test",
			Headers: H{"_hdr_": {"_catid_"}},
		},
		Response: TestResponse{
			Code: http.StatusOK,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseHeader(t *testing.T) {
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
				"X-Custom": {"custom value 123"},
			},
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseHeaderRegexp(t *testing.T) {
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
			Headers: M{
				"X-Custom": S{Regexp(`custom [a-z]+ [1-3]+`)},
			},
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseHeaderStoreVar(t *testing.T) {
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
			Headers: M{
				"X-Custom": S{StoreVar("header")},
			},
			Body: nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "custom value 123", c.r.GetVariable("header"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func JsonUseNumberUnmarshaler(data []byte, out interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(out)
}

func TestOKResponseJsonUseNumberUnmarshaler(t *testing.T) {
	// When using json UseNumber()
	// The numbers are "decoded" as json.Number (string)
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"value": 100}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:            http.StatusOK,
			BodyUnmarshaler: JsonUseNumberUnmarshaler,
			Body: M{
				"value": "100",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseRawStringBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `Hello this is plain text`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:            http.StatusOK,
			BodyUnmarshaler: RawUnmarshaler,
			Body:            "Hello this is plain text",
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseRawEmptyStringBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ``)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:            http.StatusOK,
			BodyUnmarshaler: RawUnmarshaler,
			Body:            nil,
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseRawStoreVarShortcutBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `Hello this is plain text`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:            http.StatusOK,
			BodyUnmarshaler: RawUnmarshaler,
			Body:            "$body$",
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "Hello this is plain text", c.r.GetVariable("body"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestOKResponseRawRegexpBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:            http.StatusOK,
			BodyUnmarshaler: RawUnmarshaler,
			Body:            Regexp(`^H[a-z ]+ [0-9]+$`),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKResponseRawRegexpVarsBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `Hello this is plain text 1234`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code:            http.StatusOK,
			BodyUnmarshaler: RawUnmarshaler,
			Body: RegexpVars(
				`^H[a-z ]+ ([0-9]+)$`,
				map[int]string{1: "counter"},
			),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "1234", c.r.GetVariable("counter"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestOKTestAssert(t *testing.T) {
	c := setupTest(t)

	// should not be called
	c.r.SetErrorHandler(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"ok"`)
	})

	c.r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: "ok",
		},
	})
}

func TestOKIgnoreResponseCode(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"stats": "ok"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: Any(),
			Body: M{
				"stats": "ok",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKIgnoreMapValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"stats": Any(),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKIgnoreResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Any(),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKStoreVarStringValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
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

func TestOKStoreVarNumberValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": 1580}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
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

func TestOKLoadVarStringValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "value"}`)
	})

	err := c.r.SetVariable("myvar", "value")
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
			Code: http.StatusOK,
			Body: M{
				"stats": LoadVar("myvar"),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKLoadVarNumberValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": 1580}`)
	})

	err := c.r.SetVariable("myvar", 1580)
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
			Code: http.StatusOK,
			Body: M{
				"stats": LoadVar("myvar"),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKGetVariable(t *testing.T) {
	c := setupTest(t)

	_ = c.r.SetVariable("myvar", "thevalue")

	valueStr := c.r.GetVariableString("myvar")
	if valueStr != "thevalue" {
		t.Error("value should be thevalue")
	}

	value := c.r.GetVariable("myvar")
	if value != "thevalue" {
		t.Error("value should be thevalue")
	}
}

func TestOKGetVariableNotString(t *testing.T) {
	c := setupTest(t)

	_ = c.r.SetVariable("myvar", true)

	valueStr := c.r.GetVariableString("myvar")
	if valueStr != "" {
		t.Error("value should be empty")
	}

	value := c.r.GetVariable("myvar")
	if value != true {
		t.Error("value should be true")
	}
}

func TestOKGetVariableUnknown(t *testing.T) {
	c := setupTest(t)

	valueStr := c.r.GetVariableString("myvar")
	if valueStr != "" {
		t.Error("value should be empty")
	}

	value := c.r.GetVariable("myvar")
	if value != nil {
		t.Error("value should be nil")
	}
}

func TestOKLookupVariable(t *testing.T) {
	c := setupTest(t)

	_ = c.r.SetVariable("myvar", "thevalue")

	valueStr, ok := c.r.LookupVariableString("myvar")
	if valueStr != "thevalue" {
		t.Error("value should be thevalue")
	}
	if ok != true {
		t.Error("ok should be true")
	}

	value, ok := c.r.LookupVariable("myvar")
	if value != "thevalue" {
		t.Error("value should be thevalue")
	}
	if ok != true {
		t.Error("ok should be true")
	}
}

func TestOKLookupVariableNotString(t *testing.T) {
	c := setupTest(t)

	_ = c.r.SetVariable("myvar", true)

	valueStr, ok := c.r.LookupVariableString("myvar")
	if valueStr != "" {
		t.Error("value should be empty")
	}
	if ok != false {
		t.Error("ok should be false") // value is not a string
	}

	value, ok := c.r.LookupVariable("myvar")
	if value != true {
		t.Error("value should be true")
	}
	if ok != true {
		t.Error("ok should be true") // ok as returned as interface{}
	}
}

func TestOKLookupVariableUnknown(t *testing.T) {
	c := setupTest(t)

	valueStr, ok := c.r.LookupVariableString("myvar")
	if valueStr != "" {
		t.Error("value should be empty")
	}
	if ok != false {
		t.Error("ok should be false")
	}

	value, ok := c.r.LookupVariable("myvar")
	if value != nil {
		t.Error("value should be nil")
	}
	if ok != false {
		t.Error("ok should be false")
	}
}

func TestOKRegexp(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "150 - high - end"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"stats": Regexp(`^[0-9]{3} - .* - end$`),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRegexpReplaceVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "hello world"}`)
	})

	_ = c.r.SetVariable("who", "world")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"stats": Regexp(`^hello _who_$`),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKStoreVarShortcutStringValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "high"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"stats": "$stats$",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "high", c.r.GetVariable("stats"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestOKStoreVarShortcutNumberValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": 1580}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"stats": "$stats$",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := float64(1580), c.r.GetVariable("stats"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestOKStoreVarShortcutChangedBounds(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "high"}`)
	})

	err := c.r.SetStoreShortcutBounds("(", ")")
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
			Code: http.StatusOK,
			Body: M{
				"stats": "(stats)",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "high", c.r.GetVariable("stats"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}

func TestOKLoadVarShortcutString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "value is ok !"}`)
	})

	err := c.r.SetVariable("id", "123")
	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.SetVariable("status", "ok")
	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test/_id_",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"status": "value is _status_ !",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKLoadVarShortcutInt(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "value is 100"}`)
	})

	values := []interface{}{
		int(100), int8(100), int16(100), int32(100), int64(100),
		uint(100), uint8(100), uint16(100), uint32(100), uint64(100),
	}
	for _, value := range values {
		err := c.r.SetVariable("id", value)
		if e := ExpectNil(err); e != "" {
			t.Error(e)
		}

		err = c.r.Test(TestCase{
			Request: TestRequest{
				Method: "GET",
				Path:   "/api/test/123",
				Body:   nil,
			},
			Response: TestResponse{
				Code: http.StatusOK,
				Body: M{
					"status": "value is _id_",
				},
			},
		})

		if e := ExpectNil(err); e != "" {
			t.Error(e)
		}
	}
}

func TestOKLoadVarShortcutFloat(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "value is 100.5"}`)
	})

	values := []interface{}{float32(100.5), float64(100.5)}
	for _, value := range values {
		err := c.r.SetVariable("id", value)
		if e := ExpectNil(err); e != "" {
			t.Error(e)
		}

		err = c.r.Test(TestCase{
			Request: TestRequest{
				Method: "GET",
				Path:   "/api/test/123",
				Body:   nil,
			},
			Response: TestResponse{
				Code: http.StatusOK,
				Body: M{
					"status": "value is _id_",
				},
			},
		})

		if e := ExpectNil(err); e != "" {
			t.Error(e)
		}
	}
}

func TestOKLoadVarShortcutFloatWithPrecision(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "value is 100.500"}`)
	})

	c.r.SetLoadShortcutFloatPrecision(3)

	values := []interface{}{float32(100.5), float64(100.5)}
	for _, value := range values {
		err := c.r.SetVariable("id", value)
		if e := ExpectNil(err); e != "" {
			t.Error(e)
		}

		err = c.r.Test(TestCase{
			Request: TestRequest{
				Method: "GET",
				Path:   "/api/test/123",
				Body:   nil,
			},
			Response: TestResponse{
				Code: http.StatusOK,
				Body: M{
					"status": "value is _id_",
				},
			},
		})

		if e := ExpectNil(err); e != "" {
			t.Error(e)
		}
	}
}

func TestOKLoadVarShortcutBool(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "value is true"}`)
	})

	err := c.r.SetVariable("id", true)
	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test/123",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"status": "value is _id_",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKLoadVarShortcutChangedBounds(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test/123", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "ok"}`)
	})

	err := c.r.SetVariable("id", "123")
	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.SetLoadShortcutBounds("[", "]")
	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	err = c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test/[id]",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"status": "ok",
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberDeltaExactValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberDelta(555, 0),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberDeltaLowerValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberDelta(550, 5),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberDeltaGreaterValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberDelta(560, 5),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberRange(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberRange(500, 600),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberRangeSingleElementRange(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberRange(555, 555),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberRangeLowerBound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberRange(555, 600),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKNumberRangeUpperBound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `555`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberRange(500, 555),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKTimeDeltaExactValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: TimeDelta(
				time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
				0,
			),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKTimeDeltaBeforeValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: TimeDelta(
				time.Date(2020, time.April, 11, 20, 10, 30, 0, time.UTC),
				1*time.Second,
			),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKTimeDeltaAfterValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"2020-04-11T20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: TimeDelta(
				time.Date(2020, time.April, 11, 20, 10, 31, 0, time.UTC),
				1*time.Second,
			),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKTimeDeltaDefaultFormat(t *testing.T) {
	c := setupTest(t)

	c.r.SetDefaultTimeDeltaFormat("Day 2006-01-02 Hour 15:04:05Z07:00")

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"Day 2020-04-11 Hour 20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: TimeDelta(
				time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
				0,
			),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKTimeDeltaFormat(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"Day 2020-04-11 Hour 20:10:30.123Z"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: TimeDeltaLayout(
				time.Date(2020, time.April, 11, 20, 10, 30, 123*int(time.Millisecond), time.UTC),
				0,
				"Day 2006-01-02 Hour 15:04:05Z07:00",
			),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}
}

func TestOKRegexpVars(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"The output value is: Hello and World."`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: RegexpVars(`.*: (\w+) and (\w+)\.`, map[int]string{1: "first", 2: "second"}),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "Hello", c.r.GetVariable("first"); expected != actual {
		t.Errorf("expected value %v but got %v, ", expected, actual)
	}

	if expected, actual := "World", c.r.GetVariable("second"); expected != actual {
		t.Errorf("expected value %v but got %v, ", expected, actual)
	}
}

func TestOKRegexpVarsReplaceVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "hello world"}`)
	})

	_ = c.r.SetVariable("who", "world")

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"stats": RegexpVars(`^(\w+) _who_$`, map[int]string{1: "first"}),
			},
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "hello", c.r.GetVariable("first"); expected != actual {
		t.Errorf("expected value %v but got %v, ", expected, actual)
	}

}

func TestOKRegexpVarsOnlyFullMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"--header--content--footer--"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: RegexpVars(`--header--.+--footer--`, map[int]string{0: "full"}),
		},
	})

	if e := ExpectNil(err); e != "" {
		t.Error(e)
	}

	if expected, actual := "--header--content--footer--", c.r.GetVariable("full"); expected != actual {
		t.Errorf("expected value %v but got %v", expected, actual)
	}
}
