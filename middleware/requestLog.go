package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

//https://gist.github.com/cespare/3985516
type RequestLogRecord struct {
	http.ResponseWriter

	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	elapsedTime           time.Duration
}

func (r *RequestLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	fmt.Fprintf(out, "%s - - [%s] \"%v\" %v %v %vB\n", r.ip, timeFormatted, requestLine, r.status, r.elapsedTime, r.responseBytes)
}

func (r *RequestLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

func (r *RequestLogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type RequestLogMiddleware struct {
	handler http.Handler
	out     io.Writer
}

func NewRequestLogMiddleware(handler http.Handler, out io.Writer) http.Handler {
	return &RequestLogMiddleware{
		handler: handler,
		out:     out,
	}
}

func (h *RequestLogMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &RequestLogRecord{
		ResponseWriter: rw,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.out)
}
