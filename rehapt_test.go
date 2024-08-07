package rehapt_test

import (
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
// Begin with valid cases

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
			Body: Not(10.0),
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

// And now invalid cases

func TestErrNilMarshaler(t *testing.T) {
	c := setupTest(t)

	c.r.SetMarshaler(nil)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `nil marshaler`); e != "" {
		t.Error(e)
	}
}

func TestErrNilUnmarshaler(t *testing.T) {
	c := setupTest(t)

	c.r.SetUnmarshaler(nil)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `nil unmarshaler`); e != "" {
		t.Error(e)
	}
}

func TestErrNilHTTPHandler(t *testing.T) {
	c := setupTest(t)

	c.r.SetHttpHandler(nil)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `nil HTTP handler`); e != "" {
		t.Error(e)
	}
}

func TestErrNilErrorHandler(t *testing.T) {
	server := http.NewServeMux()

	c := &testContext{
		r:      NewRehapt(nil, server),
		server: server,
	}

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"ok"`)
	})

	// The reported error on stdout here is expected
	c.r.TestAssert(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: "KO",
		},
	})

	// No easy way to check stdout, but at least we make sure the TestAssert() function
	// does not crash when errorHandler is nil
}

func TestErrMissingHTTPMethod(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `incomplete testcase. Missing HTTP method`); e != "" {
		t.Error(e)
	}
}

func TestErrInvalidHTTPMethod(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "NOT CORRECT",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `failed to build HTTP request. net/http: invalid method "NOT CORRECT"`); e != "" {
		t.Error(e)
	}
}

func TestErrMissingURLPath(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `incomplete testcase. Missing URL path`); e != "" {
		t.Error(e)
	}
}

func TestErrMarshalRequestBody(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   M{"n": json.Number(`invalid`)}, // This is refused by json.Marshal
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `failed to marshal the testcase request body. json: invalid number literal "invalid"`); e != "" {
		t.Error(e)
	}
}

func TestErrResponseCode(t *testing.T) {
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
			Code: http.StatusOK,
			Body: nil,
		},
	})

	if e := ExpectError(err, `response code does not match. Expected 200, got 401`); e != "" {
		t.Error(e)
	}
}

func TestErrTestAssertCallFailFunction(t *testing.T) {
	c := setupTest(t)

	tt := &testingT{}
	c.r.SetErrorHandler(tt)

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
			Body: "not ok",
		},
	})

	if tt.called == false {
		t.Errorf("Fail function should have been called")
	}
}

func TestErrResponseBodyType(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/string", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"ok"`)
	})
	c.server.HandleFunc("/api/int", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `1`)
	})
	c.server.HandleFunc("/api/float", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `1.0`)
	})
	c.server.HandleFunc("/api/bool", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `true`)
	})
	c.server.HandleFunc("/api/map", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"msg": "ok"}`)
	})
	c.server.HandleFunc("/api/slice", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["ok"]`)
	})

	tests := []struct {
		Path  string
		Body  interface{}
		Error string
	}{
		// Int
		{Path: "string", Body: 1, Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string"},
		{Path: "bool", Body: 1, Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got bool"},
		{Path: "map", Body: 1, Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got map"},
		{Path: "slice", Body: 1, Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got slice"},
		// Uint
		{Path: "string", Body: uint(1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string"},
		{Path: "bool", Body: uint(1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got bool"},
		{Path: "map", Body: uint(1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got map"},
		{Path: "slice", Body: uint(1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got slice"},
		// String
		{Path: "int", Body: "ok", Error: "different kinds. Expected string, got float64"}, // TODO can we force json.Unmarshal to use int ?
		{Path: "float", Body: "ok", Error: "different kinds. Expected string, got float64"},
		{Path: "bool", Body: "ok", Error: "different kinds. Expected string, got bool"},
		{Path: "map", Body: "ok", Error: "different kinds. Expected string, got map"},
		{Path: "slice", Body: "ok", Error: "different kinds. Expected string, got slice"},
		// Float32
		{Path: "string", Body: float32(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string"},
		{Path: "bool", Body: float32(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got bool"},
		{Path: "map", Body: float32(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got map"},
		{Path: "slice", Body: float32(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got slice"},
		// Float64
		{Path: "string", Body: float64(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string"},
		{Path: "bool", Body: float64(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got bool"},
		{Path: "map", Body: float64(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got map"},
		{Path: "slice", Body: float64(0.1), Error: "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got slice"},
		// Bool
		{Path: "string", Body: true, Error: "different kinds. Expected bool, got string"},
		{Path: "int", Body: true, Error: "different kinds. Expected bool, got float64"},
		{Path: "float", Body: true, Error: "different kinds. Expected bool, got float64"},
		{Path: "map", Body: true, Error: "different kinds. Expected bool, got map"},
		{Path: "slice", Body: true, Error: "different kinds. Expected bool, got slice"},
		// Map
		{Path: "string", Body: M{}, Error: "different kinds. Expected map, got string"},
		{Path: "int", Body: M{}, Error: "different kinds. Expected map, got float64"},
		{Path: "float", Body: M{}, Error: "different kinds. Expected map, got float64"},
		{Path: "bool", Body: M{}, Error: "different kinds. Expected map, got bool"},
		{Path: "slice", Body: M{}, Error: "different kinds. Expected map, got slice"},
		// Slice
		{Path: "string", Body: S{}, Error: "different kinds. Expected slice, got string"},
		{Path: "int", Body: S{}, Error: "different kinds. Expected slice, got float64"},
		{Path: "float", Body: S{}, Error: "different kinds. Expected slice, got float64"},
		{Path: "bool", Body: S{}, Error: "different kinds. Expected slice, got bool"},
		{Path: "map", Body: S{}, Error: "different kinds. Expected slice, got map"},
		// Struct
		{Path: "string", Body: struct{}{}, Error: "unhandled type struct {}"},
		{Path: "int", Body: struct{}{}, Error: "unhandled type struct {}"},
		{Path: "float", Body: struct{}{}, Error: "unhandled type struct {}"},
		{Path: "bool", Body: struct{}{}, Error: "unhandled type struct {}"},
		{Path: "slice", Body: struct{}{}, Error: "unhandled type struct {}"},
		// Unhandled
		{Path: "string", Body: complex(1, 2), Error: "unhandled type complex128"},
	}

	for _, test := range tests {
		err := c.r.Test(TestCase{
			Request: TestRequest{
				Method: "GET",
				Path:   "/api/" + test.Path,
				Body:   nil,
			},
			Response: TestResponse{
				Code: http.StatusOK,
				Body: test.Body,
			},
		})

		if e := ExpectError(err, test.Error); e != "" {
			t.Error(e)
		}
	}
}

func TestErrStringResponseBody(t *testing.T) {
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
			Body: "nok",
		},
	})

	if e := ExpectError(err, `strings does not match. Expected 'nok', got 'ok'`); e != "" {
		t.Error(e)
	}
}

func TestErrBoolResponseBody(t *testing.T) {
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
			Body: false,
		},
	})

	if e := ExpectError(err, `bools does not match. Expected false, got true`); e != "" {
		t.Error(e)
	}
}

func TestErrNotResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `10`)
	})
	c.server.HandleFunc("/api/test-str", func(w http.ResponseWriter, req *http.Request) {
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
			Body: Not(10),
		},
	})

	if e := ExpectError(err, `expected not 10, got 10`); e != "" {
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
			Body: Not("hello"),
		},
	})

	if e := ExpectError(err, `expected not hello, got hello`); e != "" {
		t.Error(e)
	}
}

func TestErrAndResponseBody(t *testing.T) {
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
			Body: And("hello", Regexp("^h...$")),
		},
	})

	if e := ExpectError(err, `regexp '^h...$' does not match 'hello'`); e != "" {
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
			Body: And("other", "unknown"),
		},
	})

	if e := ExpectError(err, `strings does not match. Expected 'other', got 'hello'`); e != "" {
		t.Error(e)
	}
}

func TestErrOrResponseBody(t *testing.T) {
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
			Body: Or("byebye", "world"),
		},
	})

	if e := ExpectError(err, `strings does not match. Expected 'byebye', got 'hello'
strings does not match. Expected 'world', got 'hello'`); e != "" {
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
			Body: Or("world", "ciao"),
		},
	})

	if e := ExpectError(err, `strings does not match. Expected 'world', got 'hello'
strings does not match. Expected 'ciao', got 'hello'`); e != "" {
		t.Error(e)
	}
}

func TestErrIntResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `100`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: 150,
		},
	})

	if e := ExpectError(err, `floats does not match. Expected 150, got 100`); e != "" {
		t.Error(e)
	}
}

func TestErrUintResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `100`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: uint(150),
		},
	})

	if e := ExpectError(err, `floats does not match. Expected 150, got 100`); e != "" {
		t.Error(e)
	}
}

func TestErrFloatResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `100.0`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: 100.5,
		},
	})

	if e := ExpectError(err, `floats does not match. Expected 100.5, got 100`); e != "" {
		t.Error(e)
	}
}

func TestErrUnmarshalResponseBody(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		// This is not valid JSON
		_, _ = fmt.Fprintf(w, `{"error": invalid...`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: Any(),
		},
	})

	if e := ExpectError(err, `cannot unmarshal response body. invalid character 'i' looking for beginning of value`); e != "" {
		t.Error(e)
	}
}

func TestErrResponseHeader(t *testing.T) {
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
				"X-Custom": {"custom value 123"},
			},
			Body: nil,
		},
	})

	if e := ExpectError(err, `response headers does not match. map element [X-Custom] does not match. slice element 0 does not match. strings does not match. Expected 'custom value 123', got 'not right value'`); e != "" {
		t.Error(e)
	}
}

func TestErrNilResponseBody(t *testing.T) {
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
			Code: http.StatusOK,
			Body: "anything",
		},
	})

	if e := ExpectError(err, `expected anything but got nil`); e != "" {
		t.Error(e)
	}
}

func TestErrRequestRawBodyInvalidType(t *testing.T) {
	c := setupTest(t)

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method:        "POST",
			Path:          "/api/test",
			BodyMarshaler: RawMarshaler,
			Body:          1,
		},
		Response: TestResponse{
			Code: http.StatusAccepted,
			Body: nil,
		},
	})

	if e := ExpectError(err, `failed to marshal the testcase request body. only string or []byte supported`); e != "" {
		t.Error(e)
	}
}

func TestErrResponseBodyExpectedNil(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"success"`)
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

	if e := ExpectError(err, `expected is nil but got success`); e != "" {
		t.Error(e)
	}
}

func TestErrSliceDifferentSize(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["A", "B"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: S{"A"},
		},
	})

	if e := ExpectError(err, `different slice sizes. Expected 1, got 2. Expected [A] got [A B]`); e != "" {
		t.Error(e)
	}
}

func TestErrSliceElementDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["A", "B"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: S{"A", "C"},
		},
	})

	if e := ExpectError(err, `slice element 1 does not match. strings does not match. Expected 'C', got 'B'`); e != "" {
		t.Error(e)
	}
}

func TestErrMapDifferentKeyType(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: map[int]string{1: "test"},
		},
	})

	if e := ExpectError(err, `different map key types. Expected int, got string`); e != "" {
		t.Error(e)
	}
}

func TestErrMapDifferentSize(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{"key": "value", "foo": "bar"},
		},
	})

	// as printed order of map is unknown, we have to expect any of the two possibilities
	e1 := ExpectError(err, `different map sizes. Expected 2, got 1. Expected map[foo:bar key:value] got map[key:value]`)
	e2 := ExpectError(err, `different map sizes. Expected 2, got 1. Expected map[key:value foo:bar] got map[key:value]`)
	if !(e1 == "" || e2 == "") {
		t.Error(e1)
	}
}

func TestErrMapKeyNotFound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{"foo": "bar"},
		},
	})

	if e := ExpectError(err, `expected key foo not found in actual map[key:value]`); e != "" {
		t.Error(e)
	}
}

func TestErrMapElementDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{"key": "bar"},
		},
	})

	if e := ExpectError(err, `map element [key] does not match. strings does not match. Expected 'bar', got 'value'`); e != "" {
		t.Error(e)
	}
}

func TestErrNumberDeltaNotNumber(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `"hi"`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberDelta(500, 10),
		},
	})

	if e := ExpectError(err, `different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string`); e != "" {
		t.Error(e)
	}
}

func TestErrNumberDeltaLowerValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `500`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberDelta(450, 49),
		},
	})

	if e := ExpectError(err, `max difference between 450 and 500 allowed is 49, but difference was 50`); e != "" {
		t.Error(e)
	}
}

func TestErrNumberDeltaGreaterValue(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `500`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: NumberDelta(550, 49),
		},
	})

	if e := ExpectError(err, `max difference between 550 and 500 allowed is 49, but difference was 50`); e != "" {
		t.Error(e)
	}
}

func TestErrTimeDeltaNotString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `1000`)
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

	if e := ExpectError(err, `different kinds. Expected string, got float64`); e != "" {
		t.Error(e)
	}
}

func TestErrTimeDeltaNotTime(t *testing.T) {
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
			Body: TimeDelta(
				time.Date(2020, time.April, 11, 20, 10, 31, 0, time.UTC),
				1*time.Second,
			),
		},
	})

	if e := ExpectError(err, `invalid time. parsing time "hello" as "2006-01-02T15:04:05Z07:00": cannot parse "hello" as "2006"`); e != "" {
		t.Error(e)
	}
}

func TestErrTimeDeltaBeforeValue(t *testing.T) {
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
				time.Date(2020, time.April, 11, 20, 10, 29, 0, time.UTC),
				1*time.Second,
			),
		},
	})

	if e := ExpectError(err, `max difference between 2020-04-11 20:10:29 +0000 UTC and 2020-04-11 20:10:30.123 +0000 UTC allowed is 1s, but difference was -1.123s`); e != "" {
		t.Error(e)
	}
}

func TestErrTimeDeltaAfterValue(t *testing.T) {
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
				time.Date(2020, time.April, 11, 20, 10, 32, 0, time.UTC),
				1*time.Second,
			),
		},
	})

	if e := ExpectError(err, `max difference between 2020-04-11 20:10:32 +0000 UTC and 2020-04-11 20:10:30.123 +0000 UTC allowed is 1s, but difference was 1.877s`); e != "" {
		t.Error(e)
	}
}

func TestErrSetVariableInvalidVarname(t *testing.T) {
	c := setupTest(t)

	err := c.r.SetVariable("my var", "value")
	if e := ExpectError(err, `invalid variable name my var`); e != "" {
		t.Error(e)
	}
}

func TestErrStoreVarInvalidVarname(t *testing.T) {
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
				"stats": StoreVar("my var"),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. invalid variable name my var`); e != "" {
		t.Error(e)
	}
}

func TestErrStoreVarInvalidBounds(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "high"}`)
	})

	err := c.r.SetStoreShortcutBounds("", ")")
	if e := ExpectError(err, `invalid prefix, cannot be empty`); e != "" {
		t.Error(e)
	}

	err = c.r.SetStoreShortcutBounds("(", "")
	if e := ExpectError(err, `invalid suffix, cannot be empty`); e != "" {
		t.Error(e)
	}
}

func TestErrLoadVarInvalidBounds(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "high"}`)
	})

	err := c.r.SetLoadShortcutBounds("", ")")
	if e := ExpectError(err, `invalid prefix, cannot be empty`); e != "" {
		t.Error(e)
	}

	err = c.r.SetLoadShortcutBounds("(", "")
	if e := ExpectError(err, `invalid suffix, cannot be empty`); e != "" {
		t.Error(e)
	}
}

func TestErrLoadVarShortcutUnknownVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "status is ok"}`)
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
				"status": "status is _unknownvar_",
			},
		},
	})

	if e := ExpectError(err, `map element [status] does not match. variable unknownvar is not defined`); e != "" {
		t.Error(e)
	}
}

func TestErrLoadVarShortcutUnknownVariableInPath(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "status is ok"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test/_unknown_",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: M{
				"status": "status is ok",
			},
		},
	})

	if e := ExpectError(err, `error while replacing variables in path. variable unknown is not defined`); e != "" {
		t.Error(e)
	}
}

func TestErrLoadVarShortcutInvalidVariableType(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status": "status is ok"}`)
	})

	err := c.r.SetVariable("var", M{"hello": "world"})
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
				"status": "status is _var_",
			},
		},
	})

	if e := ExpectError(err, `map element [status] does not match. variable var of type rehapt.M cannot be using inside string`); e != "" {
		t.Error(e)
	}
}

func TestErrUnsortedSliceDifferentSize(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["A", "B"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: UnsortedS{"A"},
		},
	})

	if e := ExpectError(err, `different slice sizes. Expected 1, got 2. Expected [A] got [A B]`); e != "" {
		t.Error(e)
	}
}

func TestErrUnsortedSliceElementNotFound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `["A", "B", "C"]`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: UnsortedS{"B", "C", "E"},
		},
	})

	if e := ExpectError(err, `expected element E at index 2 not found
actual elements at indexes [0] not found`); e != "" {
		t.Error(e)
	}
}

func TestErrPartialMapKeyNotFound(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: PartialM{"foo": "bar"},
		},
	})

	if e := ExpectError(err, `expected key foo not found`); e != "" {
		t.Error(e)
	}
}

func TestErrPartialMapElementDoesNotMatch(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "POST",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: PartialM{"key": "bar"},
		},
	})

	if e := ExpectError(err, `map element [key] does not match. strings does not match. Expected 'bar', got 'value'`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpFailParsing(t *testing.T) {
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
				"stats": Regexp(`^[0-9](3 - .* - end$`),
			},
		},
	})

	if e := ExpectError(err, "map element [stats] does not match. error parsing regexp: missing closing ): `^[0-9](3 - .* - end$`"); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpNotString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": 500}`)
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
				"stats": Regexp(`^[a-z]{3}$`),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. different kinds. Expected string, got float64`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpReplaceUnknownVariable(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"stats": "hello world"}`)
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
				"stats": Regexp(`^[a-z]+ _who_$`),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. variable who is not defined`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpDoesNotMatch(t *testing.T) {
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
				"stats": Regexp(`^[a-z]{3} - .* - end$`),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. regexp '^[a-z]{3} - .* - end$' does not match '150 - high - end'`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpVarsNotString(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `1000`)
	})

	err := c.r.Test(TestCase{
		Request: TestRequest{
			Method: "GET",
			Path:   "/api/test",
			Body:   nil,
		},
		Response: TestResponse{
			Code: http.StatusOK,
			Body: RegexpVars(`^([0-9]{3})$`, nil),
		},
	})

	if e := ExpectError(err, `different kinds. Expected string, got float64`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpVarsFailParsing(t *testing.T) {
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
				"stats": RegexpVars(`^[0-9](3 - .* - end$`, nil),
			},
		},
	})

	if e := ExpectError(err, "map element [stats] does not match. error parsing regexp: missing closing ): `^[0-9](3 - .* - end$`"); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpVarsDoesNotMatch(t *testing.T) {
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
				"stats": RegexpVars(`^[a-z]{3} - (.*) - end$`, map[int]string{1: "v1"}),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. regexp '^[a-z]{3} - (.*) - end$' does not match '150 - high - end'`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpVarsDoesInvalidVarname(t *testing.T) {
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
				"stats": RegexpVars(`^[0-9]{3} - (.*) - end$`, map[int]string{1: "v 1"}),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. invalid variable name v 1`); e != "" {
		t.Error(e)
	}
}

func TestErrRegexpVarsOverflowIndexIgnored(t *testing.T) {
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
				"stats": RegexpVars(`^[0-9]{3} - (.*) - end$`, map[int]string{2: "v1"}),
			},
		},
	})

	if e := ExpectError(err, `map element [stats] does not match. expected variable index 2 overflow regexp group count of 2`); e != "" {
		t.Error(e)
	}
}

func TestErrRawUnhandled(t *testing.T) {
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
			Body:            1234,
		},
	})

	if e := ExpectError(err, "different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got string"); e != "" {
		t.Error(e)
	}
}

func TestErrRawStringDoesNotMatch(t *testing.T) {
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
			Body:            "Hello this is plain text",
		},
	})

	if e := ExpectError(err, "strings does not match. Expected 'Hello this is plain text', got 'Hello this is plain text 1234'"); e != "" {
		t.Error(e)
	}
}

func TestErrRawRegexpFailParsing(t *testing.T) {
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
			Body:            Regexp(`^H[a-z ]+ ([0-9]+$`),
		},
	})

	if e := ExpectError(err, "error parsing regexp: missing closing ): `^H[a-z ]+ ([0-9]+$`"); e != "" {
		t.Error(e)
	}
}

func TestErrRawRegexpDoesNotMatch(t *testing.T) {
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
			Body:            Regexp(`^H[a-z ]+ [0-9]$`),
		},
	})

	if e := ExpectError(err, "regexp '^H[a-z ]+ [0-9]$' does not match 'Hello this is plain text 1234'"); e != "" {
		t.Error(e)
	}
}

func TestErrRawRegexpVarsFailParsing(t *testing.T) {
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
			Body:            RegexpVars(`^H[a-z ]+ ([0-9]+$`, map[int]string{1: "counter"}),
		},
	})

	if e := ExpectError(err, "error parsing regexp: missing closing ): `^H[a-z ]+ ([0-9]+$`"); e != "" {
		t.Error(e)
	}
}

func TestErrRawRegexpVarsDoesNotMatch(t *testing.T) {
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
			Body:            RegexpVars(`^H[a-z ]+ ([0-9])$`, map[int]string{1: "counter"}),
		},
	})

	if e := ExpectError(err, `regexp '^H[a-z ]+ ([0-9])$' does not match 'Hello this is plain text 1234'`); e != "" {
		t.Error(e)
	}
}

func TestErrRawRegexpVarsInvalidVarname(t *testing.T) {
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
			Body:            RegexpVars(`^H[a-z ]+ ([0-9]+)$`, map[int]string{1: "counter 1"}),
		},
	})

	if e := ExpectError(err, `invalid variable name counter 1`); e != "" {
		t.Error(e)
	}
}

func TestErrRawRegexpVarsOverflowIndex(t *testing.T) {
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
			Body:            RegexpVars(`^H[a-z ]+ ([0-9]+)$`, map[int]string{2: "counter"}),
		},
	})

	if e := ExpectError(err, `expected variable index 2 overflow regexp group count of 2`); e != "" {
		t.Error(e)
	}
}

func TestErrMultipleErrors(t *testing.T) {
	c := setupTest(t)

	c.server.HandleFunc("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Custom", "not right value")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"key": "value"}`)
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
			Body: M{},
		},
	})

	if e := ExpectError(err, `response code does not match. Expected 200, got 400
response headers does not match. map element [X-Custom] does not match. slice element 0 does not match. strings does not match. Expected 'custom value 123', got 'not right value'
different map sizes. Expected 0, got 1. Expected map[] got map[key:value]`); e != "" {
		t.Error(e)
	}
}
