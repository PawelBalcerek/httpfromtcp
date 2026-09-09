package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("the connection has been closed")
	}
}

func getLinesChannel(rc io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer rc.Close()

		bytes := make([]byte, 8)
		var currentLine strings.Builder
		for {
			n, err := rc.Read(bytes)
			if err != nil {
				if currentLine.Len() != 0 {
					lines <- currentLine.String()
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("an error has occurred: %s\n", err.Error())
				return
			}

			str := string(bytes[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLine.String(), parts[i])
				currentLine.Reset()
			}
			currentLine.WriteString(parts[len(parts)-1])
		}
	}()

	return lines
}
