package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Panicf("%v", err)
	}
	defer file.Close()

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(r io.Reader) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)

		bytes := make([]byte, 8)
		var currentLine strings.Builder
		for {
			n, err := r.Read(bytes)
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
