package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

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

func handleRequest(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.BadRequest,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.InternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	default:
		if _, err := fmt.Fprint(w, "All good, frfr\n"); err != nil {
			fmt.Printf("an error has occurred while writing response: %v", err)
			return &server.HandlerError{
				StatusCode: response.InternalServerError,
				Message:    "Woopsie, my bad\n",
			}
		}
		return nil
	}
}
