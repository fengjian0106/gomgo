package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go/build"

	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/stretchr/graceful"

	"github.com/fengjian0106/gomgo/context"
	"github.com/fengjian0106/gomgo/handler"
	"github.com/fengjian0106/gomgo/middleware"
)

func timeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func requestLogHandler(h http.Handler) http.Handler {
	return middleware.NewRequestLogMiddleware(h, os.Stderr)
}

/////////////////////////////
/////////////////////////////
type prefixMux []struct {
	prefix string
	h      http.Handler
}

func (pm prefixMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	founded := false
	//log.Println(r.URL.Path)
	for _, ph := range pm {
		if strings.HasPrefix(r.URL.Path, ph.prefix) {
			h = ph.h
			founded = true
			break
		}
	}

	if founded == true {
		h.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

/////////////////////////////
/////////////////////////////

func defaultBase(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

/////////////////////////////
/////////////////////////////

var (
	httpAddr  = flag.String("http", "127.0.0.1:8080", "Listen for HTTP connections on this address")
	assetsDir = flag.String("assets", filepath.Join(defaultBase("github.com/fengjian0106/gomgo"), "assets"), "Base directory for templates and static files.")
)

func main() {
	flag.Parse()

	//<0>
	//http://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
	//http://stackoverflow.com/questions/17640360/file-or-line-similar-in-golang
	log.SetFlags(log.Lshortfile)
	log.Printf("Starting server, os.Args=%s", strings.Join(os.Args, " "))

	//<1> context
	context, err := context.New()
	defer context.FreeResource()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	//<2> register api handler
	apiRouter := mux.NewRouter()

	apiRouter.Handle("/api/signin", handler.ApiHandler{context, handler.PostSigninHandler}).Methods("POST")

	apiRouter.Handle("/api/users", handler.ApiHandler{context, handler.GetUsersHandler}).Methods("GET")
	apiRouter.Handle("/api/users", handler.ApiHandler{context, handler.CreateUserHandler}).Methods("POST")
	apiRouter.Handle("/api/users/{userId}", handler.ApiHandler{context, handler.GetUserByUserIdHandler}).Methods("GET")

	apiRouter.Handle("/api/posts/{postId}", handler.ApiHandler{context, handler.GetPostByPostIdHandler}).Methods("GET")
	apiRouter.Handle("/api/posts/{postId}/comments", handler.ApiHandler{context, handler.CreateCommentForPostIdHandler}).Methods("POST")
	apiRouter.Handle("/api/users/{userId}/posts", handler.ApiHandler{context, handler.GetPostsByUserIdHandler}).Methods("GET")
	apiRouter.Handle("/api/users/{userId}/posts", handler.ApiHandler{context, handler.CreatePostHandler}).Methods("POST")

	//<3> register middleware for api handler
	// Allow 30 requests per minute
	th := throttled.RateLimit(throttled.PerMin(30), &throttled.VaryBy{RemoteAddr: true}, store.NewMemStore(1000))

	chain := alice.New(
		middleware.MakeRecoverMiddleware, //recover can not work well, TODO, FIXME
		requestLogHandler,
		middleware.MakeGzipHandler,
		middleware.MakeJsonErrorMiddleware,
		th.Throttle,
		timeoutHandler).Then(apiRouter)

	//<4> staticServer
	//this a simple static file server. If you want more control, e.g. ETag, you can use StaticServer in  http://godoc.org/github.com/golang/gddo/httputil
	fileServerHandler := http.FileServer(http.Dir(*assetsDir))

	//<5>
	/**
	if err := http.ListenAndServe(*httpAddr, chain); err != nil {
		log.Fatal(err)
	}
	*/

	log.Println("public file path is:", *assetsDir)
	//https://github.com/stretchr/graceful.git
	graceful.Run(*httpAddr, 10*time.Second, prefixMux{{"/api/", chain}, {"/public/", http.StripPrefix("/public/", fileServerHandler)}})
}
