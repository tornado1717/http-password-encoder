# http-password-encoder
An HTTP server example that listens for passwords via URL parameters and calculates their hashes.

## Installing
* Make sure the [Golang](https://golang.org/doc/install) runtime is installed on your machine.
* Clone from [github](.) or download the [source golang file](./http-password-encoder.go) directly.

## Starting up the server
0. cd into the directory with the source file.
0. run *go run http-password-encoder.go <options...>*
   * Optional parameters:
     * --port=*xxx* - specifies which port number to use (default = 8080).
     * --help - show help for command line server options.

## Accepted queries to the server
Note: These assume the server is accessible at "localhost" and was configured to use port 12345.

* ".../hash" - used to encode the provided password.  This will return the password ID that the server will use to look up the encoding after it's been processed (after 5 seconds).
  * Examples for encoding "angryMonkey" as a password:
    * From a browser: Go to "http://localhost:12345/hash?password=angryMonkey"
    * Using curl from the command line:
      * curl --data password=angryMonkey http://localhost:12345/hash
      * *or*
      * curl http://localhost:12345/hash?password=angryMonkey
* ".../hash/\<passwordID>" - this will retrieve the encoded password data once it's available (after 5 seconds of creation).
  * Examples for retrieving encoded password data for ID 1:
    * From a browser: http://localhost:12345/hash/1
    * Using curl from the command line:
      * curl http://localhost:12345/hash/1
* ".../stats" - this will return a few server statistics in JSON format. Namely, the total number of passwords encoded (including those being processed), and the total number of seconds spent handling the corresponding HTTP requests.
  * Examples for retrieving encoded password data for ID 1:
    * From a browser: http://localhost:12345/stats
    * Using curl from the command line:
      * curl http://localhost:12345/stats
* ".../shutdown" - gracefully shut down the server.  The server will immediately stop accepting new connections and will wait for all active connections to terminate before shutting down.
  * Examples for retrieving encoded password data for ID 1:
    * From a browser: http://localhost:12345/shutdown
    * Using curl from the command line:
      * curl http://localhost:12345/shutdown
