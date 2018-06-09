///////////////////////////////////////////////////////////////////////////////
// App notes
// Example usage:
//	go run cmd-line-args.go arg1 -arg2 --arg3=v3 "arg 4 with spaces"   flagArg1 -flagArg2=flagArg2_val1 --flagArg3=flagArg3_val1
//	go run cmd-line-args.go -flagArg2=flagArg2_val1 -flagArg2=flagArg2_val2 --flagArg3=flagArg3_val1 --flagArg3=flagArg3_val2
//	...
//	See: cmd-line-args__go-example-wrapper.sh
///////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"; "log"; "io"
	"io/ioutil"
	"flag"
	"net/http"
)


///////////////////////////////////////////////////////////////////////////////
// Global Constants
///////////////////////////////////////////////////////////////////////////////

const pcol string = "\x1B[32m"
const hcol string = "\x1B[33m"
const rcol string = "\x1B[0m"


///////////////////////////////////////////////////////////////////////////////
// Funcs
///////////////////////////////////////////////////////////////////////////////

func printRequest(indent string, req *http.Request) {
	fmt.Print(indent,       "req.Method='"    , req.Method    ,       "'\n")
	fmt.Print(indent,       "req.URL.*\n")
	fmt.Print(indent,       "    *.Scheme='"    , req.URL.Scheme    ,       "'\n")
	fmt.Print(indent,       "    *.Opaque='"    , req.URL.Opaque    ,       "'\n")
	fmt.Print(indent,       "    *.User='"      , req.URL.User      ,       "'\n")
	fmt.Print(indent,       "    *.Host='"      , req.URL.Host      ,       "'\n")
	fmt.Print(indent, hcol, "    *.Path='"      , req.URL.Path      , rcol, "'\n")
	fmt.Print(indent,       "    *.RawPath='"   , req.URL.RawPath   ,       "'\n")
	fmt.Print(indent,       "    *.ForceQuery='", req.URL.ForceQuery,       "'\n")
	fmt.Print(indent, hcol, "    *.RawQuery='"  , req.URL.RawQuery  , rcol, "'\n")
	fmt.Print(indent,       "    *.Fragment='"  , req.URL.Fragment  ,       "'\n")
	fmt.Print(indent,       "req.Proto='"     , req.Proto     ,       "'\n")
	fmt.Print(indent,       "req.ProtoMajor='", req.ProtoMajor,       "'\n")
	fmt.Print(indent,       "req.ProtoMinor='", req.ProtoMinor,       "'\n")
	fmt.Print(indent,       "req.Header='"    , req.Header    ,       "'\n")
	fmt.Print(indent,       "req.Body='"      , req.Body      ,       "'\n")
	{
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s    %s\n", indent, b)
	}
//	fmt.Print(indent,       "req.GetBody()='"       , req.GetBody()       ,       "'\n")
	fmt.Print(indent,       "req.ContentLength='"   , req.ContentLength   ,       "'\n")
	fmt.Print(indent,       "req.TransferEncoding='", req.TransferEncoding,       "'\n")
	fmt.Print(indent,       "req.Close='"           , req.Close           ,       "'\n")
	fmt.Print(indent,       "req.Host='"            , req.Host            ,       "'\n")
	fmt.Print(indent,       "req.Form='"            , req.Form            ,       "'\n")
	fmt.Print(indent,       "req.PostForm='"        , req.PostForm        ,       "'\n")
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
}

func handleGeneralRequest(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Message from the general request handler via io.WriteString\n")
	fmt.Fprint    (w, "Message from the general request handler via fmt.Fprint\n")
	fmt.Print(pcol, "general handler:", rcol, "\n")
	printRequest("    : ", req)
}

func hashRequested_rootOnly(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Special message from the root-only hash handler via io.WriteString\n")
	fmt.Fprint    (w, "Special message from the root-only hash handler via fmt.Fprint\n")
	fmt.Print(pcol, "root-only hash handler:", rcol, "\n")
	printRequest("    : ", req)
}

func hashRequested(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Special message from the hash handler via io.WriteString\n")
	fmt.Fprint    (w, "Special message from the hash handler via fmt.Fprint\n")
	fmt.Print(pcol, "hash handler:", rcol, "\n")
	printRequest("    : ", req)
}

func ignoreRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Print(pcol, "handler was called for a request that is being ignored..., req.URL.Path=", req.URL.Path, rcol, "\n")
}


///////////////////////////////////////////////////////////////////////////////
// Main
///////////////////////////////////////////////////////////////////////////////

func main() {
	var serverPort = flag.String("port",
		"12345",
		"Port number for the HTTP password encoder server to use",
	)
	fmt.Print("Creating server on port ", *serverPort, "\n")

	http.HandleFunc("/favicon.ico", ignoreRequest)
	http.HandleFunc("/hash/", hashRequested)
	http.HandleFunc("/hash", hashRequested_rootOnly)  // Note: if "/hash/" with a trailing slash is handled and "/hash" isn't and a client goes to "address.../hash" without the trailing slash they'll get a 301 ("Moved Permanently") error
	http.HandleFunc("/", handleGeneralRequest)
	http.HandleFunc("//", handleGeneralRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), nil))
	fmt.Println("Server created...")
} // */







///////////////////////////////////////////////////////////////////////////////
// DumpRequest() example from https://golang.org/pkg/net/http/httputil/
///////////////////////////////////////////////////////////////////////////////

/*package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%q", dump)
	}))
	defer ts.Close()

	const body = "Go is a general-purpose language designed with systems programming in mind."
	req, err := http.NewRequest("POST", ts.URL, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Host = "www.example.org"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
} // */


///////////////////////////////////////////////////////////////////////////////
// DumpResponse() example from https://golang.org/pkg/net/http/httputil/
///////////////////////////////////////////////////////////////////////////////

/*package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

func main() {
	const body = "Go is a general-purpose language designed with systems programming in mind."
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", "Wed, 19 Jul 1972 19:00:00 GMT")
		fmt.Fprintln(w, body)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q", dump)
} // */
