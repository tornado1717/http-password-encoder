///////////////////////////////////////////////////////////////////////////////
// App notes
// Example usage:
//	go run http-password-encoder.go
//	go run http-password-encoder.go --port=12345
//	...
//	See: cmd-line-args__go-example-wrapper.sh
///////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"; "log"; "io"
	"io/ioutil"
	"flag"
	"net/http"
	"net/url"
	"time"
	"os"
	"context"
)


///////////////////////////////////////////////////////////////////////////////
// Global Constants
///////////////////////////////////////////////////////////////////////////////

const (
	pcol = "\x1B[32m"
	hcol = "\x1B[33m"
	rcol = "\x1B[0m"
)


///////////////////////////////////////////////////////////////////////////////
// Funcs
///////////////////////////////////////////////////////////////////////////////

// This handles requests that weren't sent according to the project spec -
// browsers vs various curl parameters vs whatever other clients
func parseURLParams(req *http.Request, bodyData *[]byte) (url.Values, error) {
		//fmt.Println("    len(req.URL.RawQuery) =", len(req.URL.RawQuery))
		//fmt.Println("    len(*bodyData) =", len(*bodyData))

	var possibleURLParams string
	if (len(req.URL.RawQuery) > 0) {
		possibleURLParams = req.URL.RawQuery
	} else if (len(*bodyData) > 0) {
		possibleURLParams = string(*bodyData)
	}

	queryVals, err := url.ParseQuery(possibleURLParams)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return queryVals, nil
}

func printRequest(indent string, req *http.Request, bodyData *[]byte) {
	fmt.Print(indent,       "req.Method='"    , req.Method    ,       "'\n")
	fmt.Print(indent,       "req.URL.*\n")
	fmt.Print(indent,       "    *.Scheme='"        , req.URL.Scheme      ,       "'\n")
	fmt.Print(indent,       "    *.Opaque='"        , req.URL.Opaque      ,       "'\n")
	fmt.Print(indent,       "    *.User='"          , req.URL.User        ,       "'\n")
	fmt.Print(indent,       "    *.Host='"          , req.URL.Host        ,       "'\n")
	fmt.Print(indent, hcol, "    *.Path='"          , req.URL.Path        , rcol, "'\n")
	fmt.Print(indent,       "    *.RawPath='"       , req.URL.RawPath     ,       "'\n")
	fmt.Print(indent,       "    *.ForceQuery='"    , req.URL.ForceQuery  ,       "'\n")
	fmt.Print(indent, hcol, "    *.RawQuery='"      , req.URL.RawQuery    , rcol, "'\n")
	fmt.Print(indent,       "    *.Fragment='"      , req.URL.Fragment    ,       "'\n")
	fmt.Print(indent,       "req.Proto='"           , req.Proto           ,       "'\n")
	fmt.Print(indent,       "req.ProtoMajor='"      , req.ProtoMajor      ,       "'\n")
	fmt.Print(indent,       "req.ProtoMinor='"      , req.ProtoMinor      ,       "'\n")
	fmt.Print(indent,       "req.Header='"          , req.Header          ,       "'\n")
	fmt.Print(indent,       "req.Body='"            , req.Body            ,       "'\n")
	fmt.Print(indent,       "    -> '"              ,         bodyData    ,       "'\n")
	fmt.Print(indent,       "    -> '"              , string(*bodyData)   ,       "'\n")
//	fmt.Print(indent,       "req.GetBody()='"       , req.GetBody()       ,       "'\n")
	fmt.Print(indent,       "req.ContentLength='"   , req.ContentLength   ,       "'\n")
	fmt.Print(indent,       "req.TransferEncoding='", req.TransferEncoding,       "'\n")
	fmt.Print(indent,       "req.Close='"           , req.Close           ,       "'\n")
	fmt.Print(indent,       "req.Host='"            , req.Host            ,       "'\n")
	fmt.Print(indent,       "req.Host='"            , req.Host            ,       "'\n")
	fmt.Print(indent,       "<before req.ParseForm()>:\n")
	fmt.Print(indent,       "    req.Form='"            , req.Form            ,       "'\n")
	fmt.Print(indent,       "    req.PostForm='"        , req.PostForm        ,       "'\n")
	req.ParseForm()  // It seems that this only reads the URL parameters and NOT params that were thrown into the body
	fmt.Print(indent,       "<after req.ParseForm()>:\n")
	fmt.Print(indent,       "    req.Form='"            , req.Form            ,       "'\n")
	fmt.Print(indent,       "    req.PostForm='"        , req.PostForm        ,       "'\n")
	req.ParseMultipartForm(0x1000)
	fmt.Print(indent,       "<after req.ParseMultipartForm()>:\n")
	fmt.Print(indent,       "    req.Form='"            , req.Form            ,       "'\n")
	fmt.Print(indent,       "    req.PostForm='"        , req.PostForm        ,       "'\n")
	fmt.Print(indent,       "req.MultipartForm='"   , req.MultipartForm   ,       "'\n")
	fmt.Print(indent,       "req.Trailer='"         , req.Trailer         ,       "'\n")
	fmt.Print(indent,       "req.RemoteAddr='"      , req.RemoteAddr      ,       "'\n")
	fmt.Print(indent, hcol, "req.RequestURI='"      , req.RequestURI      , rcol, "'\n")
	fmt.Print(indent,       "req.TLS='"             , req.TLS             ,       "'\n")
	fmt.Print(indent,       "req.Cancel='"          , req.Cancel          ,       "'\n")
//	fmt.Print(indent,       "req.Response='"        , req.Response        ,       "'\n")
//	b, err := ioutil.ReadAll(Response.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%s", b)

	fmt.Println()
	queryVals, err := parseURLParams(req, bodyData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(indent,       "extracted query params:",       "\n")
	fmt.Print(indent,       "    error: ", err, "\n")
	fmt.Print(indent, hcol, "    queryVals: ", queryVals, rcol, "\n")
}

// Handle a few things all requests have in common:
//	* Parse any parameters given in the URL (and search the body for them if some client puts them in there for some reason).
//	* Retrieve all data from req.Body since it is an io.ReadCloser which is strictly read-once.  Data is placed into bodyData.
//	* Some general logging statements.
func processRequestCommon(w http.ResponseWriter, req *http.Request, callerNameTag string, bodyData *[]byte) {
	startTime := time.Now()

	io.WriteString(w, fmt.Sprint("Message from the ", callerNameTag, " request handler via io.WriteString\n"))
	fmt.Fprint    (w,            "Message from the ", callerNameTag, " request handler via fmt.Fprint\n")

	// Do the one time read of req.Body data
	// TODO: Should we be using (or implementing) req.GetBody() here instead?
	var err error
	*bodyData, err = ioutil.ReadAll(req.Body)
	if err != nil {
		// TODO: what if this doesn't return EOF
		log.Fatal(err)
	}
	req.Body.Close()  // the server (this) is responsible for doing this

	logTag := fmt.Sprintf("==%d:0xXXXXXXXX== %s %s",
		os.Getpid(),
		//os.GetGoRoutineID() or os.GetThreadID()  // Doesn't exist. See: https://github.com/golang/go/issues/22770
		startTime.Format("2006-01-02 15:04:05.000000000"),
		callerNameTag,
	)
	fmt.Print(pcol, callerNameTag, " handler:", rcol, "\n")
	printRequest(logTag+"    : ", req, bodyData)

	fmt.Println(logTag+"    <sleeping>")
	time.Sleep(5 * time.Second)
	fmt.Println(logTag+"    <done sleeping>")
}

func handleGeneralRequest(w http.ResponseWriter, req *http.Request) {
	var bodyData []byte
	processRequestCommon(w, req, "general", &bodyData)
}

func handleHashRequest_rootOnly(w http.ResponseWriter, req *http.Request) {
	var bodyData []byte
	processRequestCommon(w, req, "root-only hash", &bodyData)
}

func handleHashRequest(w http.ResponseWriter, req *http.Request) {
	var bodyData []byte
	processRequestCommon(w, req, "hash", &bodyData)
}

func handleIgnoredRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Print(pcol, "handler was called for a request that is being ignored..., req.URL.Path=", req.URL.Path, rcol, "\n")
}


///////////////////////////////////////////////////////////////////////////////
// Main
///////////////////////////////////////////////////////////////////////////////

func main() {
	var serverPort = flag.String("port",
		"8080",
		"Port number for the HTTP password encoder server to use",
	)
	flag.Parse()
	fmt.Print("Creating server on port ", *serverPort, "\n")

	http.HandleFunc("/favicon.ico", handleIgnoredRequest)
	http.HandleFunc("/hash/", handleHashRequest)
	http.HandleFunc("/hash", handleHashRequest_rootOnly)  // Note: if "/hash/" with a trailing slash is handled and "/hash" isn't and a client goes to "address.../hash" without the trailing slash they'll get a 301 ("Moved Permanently") error
	http.HandleFunc("/", handleGeneralRequest)  // If this doesn't happen, the default handler just returns "404 page not found"
	shutdownRequested := make(chan struct{})
	handleShutdownRequest := func(w http.ResponseWriter, req *http.Request) {
		close(shutdownRequested)
	}
	http.HandleFunc("/shutdown", handleShutdownRequest)
	http.HandleFunc("/shutdown/", handleShutdownRequest)

	server := &http.Server{Addr: fmt.Sprintf(":%s", *serverPort)}
	shutdownComplete := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
		close(shutdownComplete)
		log.Println("Server created and then shutdown...")
	}()

	<- shutdownRequested
	fmt.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}

	<- shutdownComplete  // Wait for the shutdown to finish
}
