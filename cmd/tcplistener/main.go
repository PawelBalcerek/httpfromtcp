package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PawelBalcerek/httpfromtcp/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println("a connection has been accepted")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("%v", err)
		}

		requestLine := request.RequestLine
		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			requestLine.Method,
			requestLine.RequestTarget,
			requestLine.HttpVersion,
		)

		fmt.Println("the connection has been closed")
	}
}
