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
	"flag"
	"net/http"
)


///////////////////////////////////////////////////////////////////////////////
// Global Constants
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// Funcs
///////////////////////////////////////////////////////////////////////////////

func handleGeneralRequest(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Message from the general request handler via io.WriteString\n")
	fmt.Fprint    (w, "Message from the general request handler via fmt.Fprint\n")
	fmt.Print("general handler was called..., req.URL.Path=", req.URL.Path, "\n")
}

func hashRequested_rootOnly(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Special message from the root-only hash handler via io.WriteString\n")
	fmt.Fprint    (w, "Special message from the root-only hash handler via fmt.Fprint\n")
	fmt.Print("root-only hash handler was called..., req.URL.Path=", req.URL.Path, "\n")

	handleGeneralRequest(w, req)
}

func hashRequested(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Special message from the hash handler via io.WriteString\n")
	fmt.Fprint    (w, "Special message from the hash handler via fmt.Fprint\n")
	fmt.Print("hash handler was called..., req.URL.Path=", req.URL.Path, "\n")

	handleGeneralRequest(w, req)
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

	http.HandleFunc("/hash/", hashRequested)
//	http.HandleFunc("/hash", hashRequested_rootOnly)  // Note: if "/hash/" with a trailing slash is handled and "/hash" isn't and a client goes to "address.../hash" without the trailing slash they'll get a 301 ("Moved Permanently") error
	http.HandleFunc("/", handleGeneralRequest)
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
