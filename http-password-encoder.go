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
	"context"
	"math/rand"
)


///////////////////////////////////////////////////////////////////////////////
// Globals
///////////////////////////////////////////////////////////////////////////////

const (
	pcol = "\x1B[32m"  // makes text dark green on a terminal
	hcol = "\x1B[33m"  // makes text dark yellow on a terminal
	rcol = "\x1B[0m"   // reset text colors
)


///////////////////////////////////////////////////////////////////////////////
// Funcs
///////////////////////////////////////////////////////////////////////////////

func generateGoRoutineIDTag() string {
	fgColorCode := rand.Intn(8)
	bgColorCode := rand.Intn(7)
	if fgColorCode <= bgColorCode {
		bgColorCode++
	}
	return fmt.Sprint("\x1B[", fgColorCode+30, ";", bgColorCode+40, "m##", rcol)
}

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
	log.Print(indent,       "req.Method='"    , req.Method    ,       "'\n")
	log.Print(indent,       "req.URL.*\n")
	log.Print(indent,       "    *.Scheme='"        , req.URL.Scheme      ,       "'\n")
	log.Print(indent,       "    *.Opaque='"        , req.URL.Opaque      ,       "'\n")
	log.Print(indent,       "    *.User='"          , req.URL.User        ,       "'\n")
	log.Print(indent,       "    *.Host='"          , req.URL.Host        ,       "'\n")
	log.Print(indent, hcol, "    *.Path='"          , req.URL.Path        , rcol, "'\n")
	log.Print(indent,       "    *.RawPath='"       , req.URL.RawPath     ,       "'\n")
	log.Print(indent,       "    *.ForceQuery='"    , req.URL.ForceQuery  ,       "'\n")
	log.Print(indent, hcol, "    *.RawQuery='"      , req.URL.RawQuery    , rcol, "'\n")
	log.Print(indent,       "    *.Fragment='"      , req.URL.Fragment    ,       "'\n")
	log.Print(indent,       "req.Proto='"           , req.Proto           ,       "'\n")
	log.Print(indent,       "req.ProtoMajor='"      , req.ProtoMajor      ,       "'\n")
	log.Print(indent,       "req.ProtoMinor='"      , req.ProtoMinor      ,       "'\n")
	log.Print(indent,       "req.Header='"          , req.Header          ,       "'\n")
	log.Print(indent,       "req.Body='"            , req.Body            ,       "'\n")
	log.Print(indent,       "    -> '"              ,         bodyData    ,       "'\n")
	log.Print(indent,       "    -> '"              , string(*bodyData)   ,       "'\n")
//	log.Print(indent,       "req.GetBody()='"       , req.GetBody()       ,       "'\n")
	log.Print(indent,       "req.ContentLength='"   , req.ContentLength   ,       "'\n")
	log.Print(indent,       "req.TransferEncoding='", req.TransferEncoding,       "'\n")
	log.Print(indent,       "req.Close='"           , req.Close           ,       "'\n")
	log.Print(indent,       "req.Host='"            , req.Host            ,       "'\n")
	log.Print(indent,       "req.Host='"            , req.Host            ,       "'\n")
	log.Print(indent,       "<before req.ParseForm()>:\n")
	log.Print(indent,       "    req.Form='"            , req.Form            ,       "'\n")
	log.Print(indent,       "    req.PostForm='"        , req.PostForm        ,       "'\n")
// TODO: look at req.FormValue() and req.FormPostValue()
	req.ParseForm()  // It seems that this only reads the URL parameters and NOT params that were thrown into the body
	log.Print(indent,       "<after req.ParseForm()>:\n")
	log.Print(indent,       "    req.Form='"            , req.Form            ,       "'\n")
	log.Print(indent,       "    req.PostForm='"        , req.PostForm        ,       "'\n")
	req.ParseMultipartForm(0x1000)
	log.Print(indent,       "<after req.ParseMultipartForm()>:\n")
	log.Print(indent,       "    req.Form='"            , req.Form            ,       "'\n")
	log.Print(indent,       "    req.PostForm='"        , req.PostForm        ,       "'\n")
	log.Print(indent,       "req.MultipartForm='"   , req.MultipartForm   ,       "'\n")
	log.Print(indent,       "req.Trailer='"         , req.Trailer         ,       "'\n")
	log.Print(indent,       "req.RemoteAddr='"      , req.RemoteAddr      ,       "'\n")
	log.Print(indent, hcol, "req.RequestURI='"      , req.RequestURI      , rcol, "'\n")
	log.Print(indent,       "req.TLS='"             , req.TLS             ,       "'\n")
	log.Print(indent,       "req.Cancel='"          , req.Cancel          ,       "'\n")
//	log.Print(indent,       "req.Response='"        , req.Response        ,       "'\n")
//	b, err := ioutil.ReadAll(Response.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("%s", b)

	log.Println()
	queryVals, err := parseURLParams(req, bodyData)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(indent,       "extracted query params:",       "\n")
	log.Print(indent,       "    error: ", err, "\n")
	log.Print(indent, hcol, "    queryVals: ", queryVals, rcol, "\n")
}

// Handle a few things all requests have in common:
//	* Parse any parameters given in the URL (and search the body for them if some client puts them in there for some reason).
//	* Retrieve all data from req.Body since it is an io.ReadCloser which is strictly read-once.  Data is placed into bodyData.
//	* Some general logging statements.
func processRequestCommon(w http.ResponseWriter, req *http.Request, callerNameTag, logTag string, bodyData *[]byte) {
	io.WriteString(w, fmt.Sprint("Message from the ", callerNameTag, " request handler via io.WriteString\n"))
	fmt.Fprint    (w,            "Message from the ", callerNameTag, " request handler via fmt.Fprint\n")

	// Do the one time read of req.Body data
	// TODO: Should we be using (or implementing) req.GetBody() here instead?
	var err error
	*bodyData, err = ioutil.ReadAll(req.Body)
	if err != nil {
		// TODO: what if this doesn't return EOF
		log.Fatal(fmt.Sprint(logTag, err))
	}
	req.Body.Close()  // the server (this) is responsible for doing this

	log.Print(logTag, pcol, " ", callerNameTag, " handler:", rcol, "\n")
	printRequest(logTag + "    : ", req, bodyData)

/*	log.Println(logTag, "    <sleeping>")
	time.Sleep(5 * time.Second)
	log.Println(logTag, "    <done sleeping>") */
}

func handleGeneralRequest(w http.ResponseWriter, req *http.Request) {
	callerNameTag := "general"
	logTag := generateGoRoutineIDTag() + " " + callerNameTag
	var bodyData []byte
	processRequestCommon(w, req, "general", logTag, &bodyData)
}

func handleHashRequest_rootOnly(w http.ResponseWriter, req *http.Request) {
	callerNameTag := "root-only hash"
	logTag := generateGoRoutineIDTag() + " " + callerNameTag

	var bodyData []byte
	processRequestCommon(w, req, callerNameTag, logTag, &bodyData)
}

func handleHashRequest(w http.ResponseWriter, req *http.Request) {
	callerNameTag := "hash"
	logTag := generateGoRoutineIDTag() + " " + callerNameTag

	var bodyData []byte
	processRequestCommon(w, req, callerNameTag, logTag, &bodyData)
}

func handleIgnoredRequest(w http.ResponseWriter, req *http.Request) {
	log.Print(pcol, "handler was called for a request that is being ignored..., req.URL.Path=", req.URL.Path, rcol, "\n")
}


///////////////////////////////////////////////////////////////////////////////
// Main
///////////////////////////////////////////////////////////////////////////////

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	var serverPort = flag.String("port",
		"8080",
		"Port number for the HTTP password encoder server to use",
	)
	flag.Parse()
	log.Print("Creating server on port ", *serverPort, "\n")

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
	log.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}

	<- shutdownComplete  // Wait for the shutdown to finish
}
