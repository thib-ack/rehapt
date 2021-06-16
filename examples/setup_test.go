package examples

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/thib-ack/rehapt"
)

// We create our Rehapt instance with little customization for tests
func setupRehapt(t *testing.T) *rehapt.Rehapt {
	r := rehapt.NewRehapt(t, httpServer())
	// Customize a bit for our tests
	r.SetDefaultHeader("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
	return r
}

// This is our example HTTP test server.

func httpServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/user", func(w http.ResponseWriter, req *http.Request) {
		// This API support only GET
		if req.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, `{"error": "not found"}`)
			return
		}

		// This API requires Auth
		if u, p, ok := req.BasicAuth(); ok == false || u != "Aladdin" || p != "open sesame" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, `{"error": "unauthorized"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"id": "55", "name": "John", "age": 51, "pets": [{"id": "123", "type": "cat", "name": "Pepper the cat"}], "weddingdate": "2019-06-22T16:00:10.123Z"}`)
	})

	mux.HandleFunc("/api/cat/123", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Pet-Type", "Cat")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"id": "123", "name": "Pepper the cat", "age": 3, "owner": {"id": "55", "name": "John", "age": 51}, "toys": ["ball", "plastic mouse"]}`)
	})

	return mux
}
