package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/fengjian0106/gomgo/handler"
)

type JsonErrorMiddleware struct {
	handler http.Handler
}

func MakeJsonErrorMiddleware(h http.Handler) http.Handler {
	return &JsonErrorMiddleware{h}
}

func (m *JsonErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	// passing a ResponseRecorder instead of the original RW
	m.handler.ServeHTTP(rec, r)
	// after this finishes, we have the response recorded
	// and can modify it before copying it to the original RW

	if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") &&
		((rec.Code >= http.StatusBadRequest && rec.Code <= http.StatusTeapot) ||
			(rec.Code >= http.StatusInternalServerError && rec.Code <= http.StatusHTTPVersionNotSupported)) {
		//response is not json AND status code >= 400/500, change the response info to apierror

		w.WriteHeader(rec.Code)

		// we copy the original headers first
		for k, v := range rec.Header() {
			if k == "Content-Type" {
				// and set an additional one
				log.Println("json error middleware, set header to json")
				w.Header().Set("Content-Type", "application/json")
			} else {
				w.Header()[k] = v
			}
		}

		// chagne the body to json
		errJson := fmt.Sprintf("{code: %d, err: %s}\n", handler.ApiErrorUnknow, rec.Body.Bytes())
		w.Write([]byte(errJson))

	} else {
		// we copy the original headers first
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		// then write out the original body
		w.Write(rec.Body.Bytes())
	}
}
