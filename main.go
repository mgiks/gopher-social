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
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(c, "error reading command: %v\n", err)
		return
	}

	parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(parts) != 2 {
		fmt.Fprintf(c, "invalid command format: expected format 'COMMAND RESOURCE'\n")
		return
	}

	command := parts[0]
	resource := parts[1]
	log.Printf("Received command: %s %s\n", command, resource)

	switch command {
	case "GET":
		handleGet(c, resource)
	default:
		fmt.Fprintf(c, "Unknown commmand: %s\n", command)
	}
}

func handleGet(c net.Conn, resource string) {
	fmt.Fprintf(c, "GET command received for resource: %s\n", resource)
}
