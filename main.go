package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// Handler error
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	// Creates a new reader from the connection
	reader := bufio.NewReader(c)

	// Read the command line from the client
	line, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprintf(c, "error reading command: %v\n", err)
		return
	}

	parts := strings.SplitN(strings.TrimSpace(string(line)), " ", 2)
	if len(parts) != 2 {
		fmt.Fprintln(c, "invalid command format: expected format 'COMMAND RESOURCE'")
		return
	}
}
