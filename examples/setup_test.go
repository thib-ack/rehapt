package examples

import (
	"fmt"
	"net/http"
	"rehapt"
	"testing"
)

// We create our Rehapt instance with little customization for tests
func setupRehapt(t *testing.T) *rehapt.Rehapt {
	r := rehapt.NewRehapt(httpServer())
	// Customize a bit for our tests
	r.SetDefaultHeader("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
	r.SetFail(func(err error) {
		// Unfortunately all errors will be reported from this line here and not the real test which triggered it.
		// This can be solved by using the runtime.Callers() method
		// or simply use great libs like github.com/stretchr/testify/assert
		// where a simple assert.Fail(t, err) will fix the issue
		t.Error(err)
	})
	return r
}

// This is our example HTTP test server.

func httpServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/user", func(w http.ResponseWriter, req *http.Request) {
		// This API support only GET
		if req.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": "not found"}`)
			return
		}

		// This API requires Auth
		if u, p, ok := req.BasicAuth(); ok == false || u != "Aladdin" || p != "open sesame" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "unauthorized"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"id": "55", "name": "John", "age": 51, "pets": [{"id": "123", "type": "cat", "name": "Pepper the cat"}], "weddingdate": "2019-06-22T16:00:10.123Z"}`)
	})

	mux.HandleFunc("/api/cat/123", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Pet-Type", "Cat")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"id": "123", "name": "Pepper the cat", "age": 3, "owner": {"id": "55", "name": "John", "age": 51}, "toys": ["ball", "plastic mouse"]}`)
	})

	return mux
}
