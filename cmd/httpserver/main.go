package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/PawelBalcerek/httpfromtcp/internal/headers"
	"github.com/PawelBalcerek/httpfromtcp/internal/request"
	"github.com/PawelBalcerek/httpfromtcp/internal/response"
	"github.com/PawelBalcerek/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handleRequest)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port: ", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

const (
	HTTP_BIN_PREFIX = "/httpbin/"

	YOUR_PROBLEM_PAYLOAD = `
	<html>
  		<head>
    		<title>400 Bad Request</title>
  		</head>
  		<body>
    		<h1>Bad Request</h1>
    		<p>Your request honestly kinda sucked.</p>
  		</body>
	</html>
	`

	MY_PROBLEM_PAYLOAD = `
	<html>
	  	<head>
	    	<title>500 Internal Server Error</title>
		</head>
		<body>
	    	<h1>Internal Server Error</h1>
	    	<p>Okay, you know what? This one is on me.</p>
	  	</body>
	</html>
	`

	OK_PAYLOAD = `
	<html>
		<head>
	    	<title>200 OK</title>
	  	</head>
	  	<body>
	    	<h1>Success!</h1>
	    	<p>Your request was an absolute banger.</p>
	  	</body>
	</html>
	`
)

func handleRequest(w *response.Writer, req *request.Request) {
	requestTarget := req.RequestLine.RequestTarget
	if strings.HasPrefix(requestTarget, HTTP_BIN_PREFIX) {
		proxyHandler(w, requestTarget)
		return
	}

	var p string
	switch requestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.BadRequest)
		p = YOUR_PROBLEM_PAYLOAD
	case "/myproblem":
		w.WriteStatusLine(response.InternalServerError)
		p = MY_PROBLEM_PAYLOAD
	default:
		w.WriteStatusLine(response.Ok)
		p = OK_PAYLOAD
	}

	h := headers.NewHeaders()
	h.SetConnection("close")
	h.SetContentType("text/html")
	p = strings.TrimSpace(p)
	h.SetContentLength(len(p))
	w.WriteHeaders(h)
	w.WriteBody(p)
}

func proxyHandler(w *response.Writer, requestTarget string) {
	resp, err := http.Get(fmt.Sprintf("https://httpbingo.org/%s", strings.TrimPrefix(requestTarget, HTTP_BIN_PREFIX)))
	if err != nil {
		fmt.Printf("an error has occurred while calling httpbingo: %v\n", err)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.Ok)

	h := headers.NewHeaders()
	h.SetConnection("close")
	h.SetContentType("application/json")
	h.SetHeader("Transfer-Encoding", "chunked")
	h.SetHeader("Trailer", "X-Content-SHA256")
	h.SetHeader("Trailer", "X-Content-Length")
	w.WriteHeaders(h)

	buffer := make([]byte, 1024)
	wroteBytes := 0
	hasher := sha256.New()
	for {
		readBytes, err := resp.Body.Read(buffer)

		if readBytes > 0 {
			n, err := w.WriteChunkedBody(buffer[:readBytes])
			if err != nil {
				fmt.Printf("an error has occurred while writing chunked response from httpbingo to writer: %v\n", err)
				return
			}

			if _, err = hasher.Write(buffer[:readBytes]); err != nil {
				fmt.Printf("an error has occurred while writing chunked response from httpbingo to hasher: %v\n", err)
				return
			}

			wroteBytes += n
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				t := headers.NewHeaders()
				t.SetHeader("x-content-sha256", hex.EncodeToString(hasher.Sum(nil)))
				t.SetHeader("x-content-length", strconv.Itoa(wroteBytes))
				w.WriteTrailers(t)
				return
			}
			fmt.Printf("an error has occurred while reading response from httpbingo: %v\n", err)
			return
		}
	}
}
