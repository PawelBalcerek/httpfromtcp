package main

import (
	"log"
	"os"
	"os/signal"
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
	var p string
	switch req.RequestLine.RequestTarget {
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
