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
	"net/url"
//	"strings"
//	"bytes"  // needed to implement getIOReaderLen()
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

/*
			func getIOReaderLen(rc io.ReadCloser) (int64) {
				// Modified from https://stackoverflow.com/questions/39064343/how-to-get-the-size-of-an-io-reader-object
				//	"Since io.Reader interface doesn't know anything about size or length of underlying data..."
				//		Ex: a Reader's data could be coming from a stream -> unknown final length
				//	I think they chose io.Copy() over io.Read...() since it uses the correct option between "src.WriteTo(dst)" and "dst.ReadFrom(src)"
				// Also see:
				//	https://stackoverflow.com/questions/30910487/why-an-io-reader-after-read-it-became-empty-in-golang

				tmpBuf := &bytes.Buffer{}  // Need a type (like bytes.Buffer) that can be used as an io.Writer (something that implements Write())
				nRead, err := io.Copy(tmpBuf, rc)
				if err != nil {
					fmt.Println(err)
				}

				return nRead // effectively len(rc)
			}
			//func len()
			/*func (rc *io.Reader) Len() int {
				return 0
			} * /
*/

// This handles requests that weren't sent according to the project spec -
// browsers vs various curl parameters vs whatever other clients
func parseURLParams(req *http.Request, bodyData *[]byte) (url.Values, error) {
		fmt.Println("    len(req.URL.RawQuery) =", len(req.URL.RawQuery))
//		fmt.Println("    getIOReaderLen(req.Body) =", getIOReaderLen(req.Body))
		fmt.Println("    len(*bodyData) =", len(*bodyData))

	var queryVals url.Values
	var err error

	if (len(req.URL.RawQuery) > 0) {
		var err error
		queryVals, err = url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	//} else if ((&req.Body).Len() > 0) {
	//} else if (len(req.Body) > 0) {
	//} else if (getIOReaderLen(req.Body) > 0) {
	} else if (len(*bodyData) > 0) {
		/*bodyData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		} */

		queryVals, err = url.ParseQuery(string(*bodyData))
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	//a := "abc"
	//fmt.Printf("%T, %v, %s, %d", a, a, a, a.Len())

	return queryVals, nil
}

func printRequest(indent string, req *http.Request, bodyData *[]byte) {
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
	fmt.Print(indent,       "req.Body='"      , req.Body      ,       "'\n")
/*	{  // req.Body is an io.ReadCloser which is strictly read-once so special handling is required
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s    '%s'\n", indent, b)

		b, err = ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s    '%s'\n", indent, b)

		b, err = ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s    '%s'\n", indent, b)
	} */
	fmt.Print(indent,       "    -> '"              ,         bodyData    ,       "'\n")
	fmt.Print(indent,       "    -> '"              , string(*bodyData)   ,       "'\n")
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

	fmt.Println()
	queryVals, err := parseURLParams(req, bodyData)
	fmt.Print(indent,       "extracted query params:",       "\n")
	fmt.Print(indent,       "    error: ", err, "\n")
	fmt.Print(indent, hcol, "    queryVals: ", queryVals, rcol, "\n")
}

// Handle a few things all requests have in common:
//	* Parse any parameters given in the URL (and search the body for them if some client puts them in there for some reason).
//	* Retrieve all data from req.Body since it is an io.ReadCloser which is strictly read-once.  Data is placed into bodyData.
//	* Some general logging statements.
func processRequestCommon(w http.ResponseWriter, req *http.Request, callerNameTag string, bodyData *[]byte) () {
	io.WriteString(w, fmt.Sprint("Message from the ", callerNameTag, " request handler via io.WriteString\n"))
	fmt.Fprint    (w,            "Message from the ", callerNameTag, " request handler via fmt.Fprint\n")

	// Do the one time read of req.Body data
	var err error
	*bodyData, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	req.Body.Close()  // the server (this) is responsible for doing this
	//fmt.Printf("%s    '%s'\n", indent, bodyData)

	fmt.Print(pcol, callerNameTag, " handler:", rcol, "\n")
	printRequest("    : ", req, bodyData)
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
	processRequestCommon(w, req, "hash handler", &bodyData)
}

func handleIgnoredRequest(w http.ResponseWriter, req *http.Request) {
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

	http.HandleFunc("/favicon.ico", handleIgnoredRequest)
	http.HandleFunc("/hash/", handleHashRequest)
	http.HandleFunc("/hash", handleHashRequest_rootOnly)  // Note: if "/hash/" with a trailing slash is handled and "/hash" isn't and a client goes to "address.../hash" without the trailing slash they'll get a 301 ("Moved Permanently") error
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
