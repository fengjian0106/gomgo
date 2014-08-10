package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type RecoverMiddleware struct {
	handler http.Handler
}

func MakeRecoverMiddleware(h http.Handler) http.Handler {
	return &RecoverMiddleware{h}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func testPanic() {
	panic("!!test panics")
}

//https://github.com/go-martini/martini/blob/master/recovery.go
func (m *RecoverMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Catch and log any panics
	//log.Println("RecoverMiddleware::ServeHTTP")
	defer func() {
		if err := recover(); err != nil {
			stack := stack(3)
			log.Printf("------!!!!!!  PANIC: %s\n%s", err, stack)
		}
	}()

	//TODO, when panics, if need to return 500 to http client??
	m.handler.ServeHTTP(w, r)

	/**
	if code is
	testPanic()
	the defer can recover the test panic


	but if code is
	go testPanic()
	the defer CAN NOT recover the test panic


	So, the defer can not recover the panic from call  m.handler.ServeHTTP(w, r)
	I don't know why, maybe I am wrong
	If I am ringt, there is no need to use this recovermiddleware
	//FIXME
	*/
}
