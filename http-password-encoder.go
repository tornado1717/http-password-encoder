///////////////////////////////////////////////////////////////////////////////
// App notes
// Example usage:
//	go run http-password-encoder.go
//	go run http-password-encoder.go --port=12345
///////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"; "log"; "io"
	"io/ioutil"
	"flag"
	"net/http"; "net/url"
	"context"; "sync"
	"math/rand"
	"time"
	"crypto/sha512"; "encoding/base64"
	"regexp"; "strconv"
	"encoding/json"
)


///////////////////////////////////////////////////////////////////////////////
// Globals
///////////////////////////////////////////////////////////////////////////////

const (
	pcol = "\x1B[32m"  // makes text dark green on a terminal
	hcol = "\x1B[33m"  // makes text dark yellow on a terminal
	rcol = "\x1B[0m"   // reset text colors
)

var passwordDataMutex sync.RWMutex
type passwordDataElement struct {
	requestedTime time.Time
	encodedPasswordHash string
}
var passwordData []passwordDataElement

var serverStats struct {
	sync.RWMutex
	newHashNumRequests int
	newHashNSecs int64
} = struct {
	sync.RWMutex
	newHashNumRequests int
	newHashNSecs int64
}{
	newHashNumRequests: 0,
	newHashNSecs: 0,
}


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

	// Note:
	//	In the case of browsers (Chrome, Firefox) or "curl http://...<ADDRESS>...?paramName=paramValue",
	//		the request will be a "GET" request and the parameters will be stored in the URL.RawQuery field
	//	In the case of "curl --data paramName=paramValue http://...<ADDRESS>...",
	//		the request will be a "POST" request and the parameters will be stored in the request body
	var possibleURLParams string
	if (len(req.URL.RawQuery) > 0) {
		possibleURLParams = req.URL.RawQuery
	} else if (len(*bodyData) > 0) {
		possibleURLParams = string(*bodyData)
	}

	queryVals, err := url.ParseQuery(possibleURLParams)
	if err != nil {
		log.Print(err)
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
//		log.Print(err)
//		return
//	}
//	log.Printf("%s", b)

	log.Println(indent, "-----")
	queryVals, err := parseURLParams(req, bodyData)
	if err != nil {
		log.Print(err)
		return
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
		log.Print(logTag, err)
		return
	}
	req.Body.Close()  // the server (this) is responsible for doing this

	log.Print(logTag, pcol, " ", callerNameTag, " handler:", rcol, "\n")
	printRequest(logTag + "    : ", req, bodyData)

/*	log.Println(logTag, "    <sleeping>")
	time.Sleep(5 * time.Second)
	log.Println(logTag, "    <done sleeping>") */
}

// Handle requests that weren't sent to one of the pre-defined end points
func handleGeneralRequest(w http.ResponseWriter, req *http.Request) {
	callerNameTag := "general"
	logTag := generateGoRoutineIDTag() + " " + callerNameTag
	var bodyData []byte
	processRequestCommon(w, req, "general", logTag, &bodyData)
	http.NotFound(w, req)
	fmt.Fprint(w,
		"The only endpoints available on this server are:\n" +
		"    .../hash - used to encode the provided password\n" +
		"    .../hash/<passwordID> - this will retrieve the encoded password data once it is available\n" +
		"    .../stats - this will return a few server statistics in JSON format\n" +
		"    .../shutdown - gracefully shut down the server.  The server will immediately stop accepting new connections and will wait for all active connections to terminate before shutting down\n")
}

// Store and encode a new password if provided
func handleHashRequest_rootOnly(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()

	callerNameTag := "root-only hash"
	logTag := generateGoRoutineIDTag() + " " + callerNameTag

	var bodyData []byte
	processRequestCommon(w, req, callerNameTag, logTag, &bodyData)
	queryVals, err := parseURLParams(req, &bodyData)
	if err != nil {
		log.Print(logTag, err)
		return
	}

	paramVals, paramExists := queryVals["password"]
	if (paramExists) {
		pwPlainText := paramVals[0]
		hashVal := sha512.Sum512([]byte(pwPlainText))

		passwordDataMutex.Lock()
		passwordData = append(passwordData, passwordDataElement{
			requestedTime: startTime,
			encodedPasswordHash: base64.StdEncoding.EncodeToString(hashVal[:]),
		})
		passID := len(passwordData)
		passwordDataMutex.Unlock()
fmt.Println(len(passwordData), passwordData[len(passwordData)-1].encodedPasswordHash)
		fmt.Fprint(w, passID, "\n")
	} else {
		fmt.Fprint(w, "Error: no password given\n")
	}

	serverStats.Lock()
	serverStats.newHashNumRequests++
	serverStats.newHashNSecs += time.Now().Sub(startTime).Nanoseconds()
	serverStats.Unlock()
}

// Returns a request handler function that will return the encoded password corresponding with the ID provided (if available)
// If no ID was provided, store new password (also if provided)
func makeFunc_handleHashRequest() func
	(w http.ResponseWriter, req *http.Request,
) {
	compiledNewHashRegex := regexp.MustCompile("/hash/*$")
	compiledHashIDRegex  := regexp.MustCompile("/hash/([0-9]+)/*$")

	return func (w http.ResponseWriter, req *http.Request) {
		if compiledNewHashRegex.MatchString(req.RequestURI) {
			handleHashRequest_rootOnly(w, req)
			return
		} else if matchedStrings := compiledHashIDRegex.FindStringSubmatch(req.RequestURI); len(matchedStrings) == 2 {
			callerNameTag := "hash"
			logTag := generateGoRoutineIDTag() + " " + callerNameTag

			i, err := strconv.Atoi(matchedStrings[1])
			if err != nil {
				http.NotFound(w, req)  // somehow this didn't translate (the previous regex should have made this impossible) - send 404 response
				return
			}

			var encodedPasswordHash string
			encodedPasswordHashAvailable := false
			passwordDataMutex.RLock()
			if (1 <= i) && (i <= len(passwordData)) && (time.Since(passwordData[i-1].requestedTime).Seconds() >= 5) {
				encodedPasswordHashAvailable = true
				encodedPasswordHash = passwordData[i-1].encodedPasswordHash
			}
			passwordDataMutex.RUnlock()

			if encodedPasswordHashAvailable {
				var bodyData []byte
				processRequestCommon(w, req, callerNameTag, logTag, &bodyData)
				io.WriteString(w, encodedPasswordHash + "\n")
				return
			} else {
				http.NotFound(w, req)  // no pattern matched (no corresponding hash ID available yet) - send 404 response
				return
			}
		} else {
			http.NotFound(w, req)  // no pattern matched (an invalid ID was given) - send 404 response
			return
		}
	}
}

// Handle requests that are to be specifically ignored
func handleIgnoredRequest(w http.ResponseWriter, req *http.Request) {
	log.Print(pcol, "handler was called for a request that is being ignored..., req.URL.Path=", req.URL.Path, rcol, "\n")
}

// Returns some server statistics
func handleStatsRequest(w http.ResponseWriter, req *http.Request) {
	callerNameTag := "stats"
	logTag := generateGoRoutineIDTag() + " " + callerNameTag
	var bodyData []byte
	processRequestCommon(w, req, "stats", logTag, &bodyData)

	serverStats.RLock()
	stats := struct {
		Total int  "total"  // Need to give this a field lable to get "total" to have a lower case name in the JSON output
		Average float32  "average"  // Ditto for "average"
	}{
		Total: serverStats.newHashNumRequests,
		Average: float32(serverStats.newHashNSecs) / (float32(time.Second) * float32(serverStats.newHashNumRequests)),
	}
	serverStats.RUnlock()

	encodedStats, err := json.Marshal(stats)
fmt.Println("encodedStats:", encodedStats)
	if err == nil {
		io.WriteString(w, string(encodedStats) + "\n")
	}
}


///////////////////////////////////////////////////////////////////////////////
// Main
///////////////////////////////////////////////////////////////////////////////

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)


	////////////////////
	// Parse command line parameters
	////////////////////

	var showHelp = flag.Bool("help",
		false,
		"Show help",
	)
	var serverPort = flag.Int("port",
		8080,
		"Port number for the HTTP password encoder server to use",
	)
	flag.Parse()
	if *showHelp {
		flag.PrintDefaults()
		return
	}
	if (0 > *serverPort) || (*serverPort > 65535) {
		fmt.Println("Invalid port number - must be in the range [0, 65535]")
		return
	}
	log.Print("Creating server on port ", *serverPort, "\n")


	////////////////////
	// Initialize settings and handlers for HTTP server
	////////////////////

	// Note:
	//	If "/hash" without a trailing slash is handled:
	//		* and "/hash/" isn't and a client goes to:
	//			"address.../hash/" or
	//			"address.../hash/any" or
	//			"address.../hash/any/" or
	//			"address.../hash/any/thingElse" or
	//			"address.../hash/any/thingElse/"
	//			* and "/" is handled
	//				the request will be handled by the general "/" handler
	//			* and "/" isn't handled
	//				they'll get a 404 ("page not found") error sent back to them by this server
	//	If "/hash/" with a trailing slash is handled:
	//		* and "/hash" isn't
	//			* whether or not "/" is and a client goes to:
	//				*
	//					"address.../hash" without the trailing slash
	//					they'll get a 301 ("Moved Permanently") error sent back to them by this server
	//				*
	//					"address.../hash/" or
	//					"address.../hash/any" or
	//					"address.../hash/any/" or
	//					"address.../hash/any/thingElse" or
	//					"address.../hash/any/thingElse/"
	//					the request will be handled by the "/hash/" handler
	//			* and "/" isn't and a client goes to:
	//				<almost anything except "/hash/..." or other specifically handled URLs>
	//				they'll get a 404 ("page not found") error sent back to them
	//
	//	If "/" is handled:
	//			* and "/" is and a client goes to:
	//				"address..." or      <- this path will be implicitly changed to have a trailing slash added to it by curl or Chrome or Firefox or ...
	//				"address.../" or
	//				"address.../somethingElse" or
	//				"address.../somethingElseWithSlash/" or
	//				"address.../whatever/else/with/or/without/trailing/slash"
	//				the request will be handled by the general "/" handler
	http.HandleFunc("/favicon.ico", handleIgnoredRequest)
	http.HandleFunc("/hash/", makeFunc_handleHashRequest())
	http.HandleFunc("/hash", handleHashRequest_rootOnly)
	http.HandleFunc("/stats/", handleStatsRequest)
	http.HandleFunc("/stats", handleStatsRequest)
	http.HandleFunc("/", handleGeneralRequest)  // If this doesn't happen, the default handler just returns "404 page not found"
	shutdownRequested := make(chan string)
	handleShutdownRequest := func(w http.ResponseWriter, req *http.Request) {
		close(shutdownRequested)
	}
	http.HandleFunc("/shutdown", handleShutdownRequest)
	http.HandleFunc("/shutdown/", handleShutdownRequest)

	server := &http.Server{Addr: fmt.Sprintf(":%d", *serverPort)}

	// Tell the server not to reuse underlying TCP connections - makes server shutdown quick
	//	since none of the requests leave connections around that need to time out and this
	//	application probably doesn't benefit much from reuse of connections anyways due to its
	//	small data usage per request
	// Note/TODO: this can cause "TIME_WAIT" connections to be left around - is this a problem?
	//	See: http://www.serverframework.com/asynchronousevents/2011/01/time-wait-and-its-design-implications-for-protocols-and-scalable-servers.html
	server.SetKeepAlivesEnabled(false)


	////////////////////
	// Run HTTP server and until a shutdown request comes
	////////////////////

	const serverInitError = "Initialization error"
	shutdownComplete := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("HTTP server ListenAndServe: %v", err)
			if err != http.ErrServerClosed {
				shutdownRequested <- serverInitError
				close(shutdownRequested)
			}
		}
		close(shutdownComplete)
		log.Println("Server created and then shutdown...")
	}()

	shutdownRequestCause := <- shutdownRequested
	if shutdownRequestCause != serverInitError {
		log.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

		<- shutdownComplete  // Wait for the shutdown to finish
	}
}
