package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"code.google.com/p/go.net/context"

	"github.com/fengjian0106/gomgo/appcontext"
	"github.com/fengjian0106/gomgo/ctxutil"
)

func GoogleSearchHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, req *http.Request) (int, error) {
	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(req.FormValue("timeout"))
	if err == nil {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires.
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	// Check the search query.
	query := req.FormValue("q")
	if query == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorGoogleSearchQueryWordNotFound, errors.New("no query")}
	}

	// Store the user IP in ctx for use by code in other packages.
	userIP, err := ctxutil.IpFromRequest(req)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorGoogleSearchErr, err}
	}
	ctx = ctxutil.NewContextWithIp(ctx, userIP)

	// Run the Google search and print the results.
	start := time.Now()
	results, err := googleSearch(ctx, query)
	elapsed := time.Since(start)
	//log.Println(elapsed)
	if err != nil {
		return http.StatusInternalServerError, &ApiError{ApiErrorGoogleSearchErr, err}
	}

	type ReturnValue struct {
		Data           *Results `json:"data"`
		ElapsedSeconds float64  `json:"elapsedSeconds"`
	}
	rv := ReturnValue{&results, elapsed.Seconds()}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&rv)

	return http.StatusOK, nil
}

////////////////////
// Do Google searches using the Google Web Search API.
// See https://developers.google.com/web-search/docs/
////////////////////
// Results is an ordered list of search results.
type Results []Result

// A Result contains the title and URL of a search result.
type Result struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Search sends query to Google search and returns the results.
func googleSearch(ctx context.Context, query string) (Results, error) {
	// Prepare the Google Search API request.
	req, err := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)

	// If ctx is carrying the user IP address, forward it to the server.
	// Google APIs use the user IP to distinguish server-initiated requests
	// from end-user requests.
	if userIP, ok := ctxutil.IpFromContext(ctx); ok {
		q.Set("userip", userIP.String())
	}
	req.URL.RawQuery = q.Encode()

	// Issue the HTTP request and handle the response. The httpDo function
	// cancels the request if ctx.Done is closed.
	var results Results
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Parse the JSON search result.
		// https://developers.google.com/web-search/docs/#fonje
		var data struct {
			ResponseData struct {
				Results []struct {
					TitleNoFormatting string
					URL               string
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		for _, res := range data.ResponseData.Results {
			results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
		}
		return nil
	})
	// httpDo waits for the closure we provided to return, so it's safe to
	// read results here.
	return results, err
}

// httpDo issues the HTTP request and calls f with the response. If ctx.Done is
// closed while the request or f is running, httpDo cancels the request, waits
// for f to exit, and returns ctx.Err. Otherwise, httpDo returns f's error.
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
