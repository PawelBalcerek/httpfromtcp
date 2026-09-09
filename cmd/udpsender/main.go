package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("%v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		str, err := r.ReadString('\n')
		if err != nil {
			log.Printf("an error during read has occurred: %v", err)
		}
		if _, err = conn.Write([]byte(str)); err != nil {
			log.Printf("an error during write has occurred: %v", err)
		}
	}
}
