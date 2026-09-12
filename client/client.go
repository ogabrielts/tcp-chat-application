package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to connect to server. %s\n", err)
	}
	defer conn.Close()

	fmt.Print("Enter username: ")

	// Receive messages from the server
	go func () {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil { 
				if err == io.EOF {
					log.Fatalf("%s\n", msg)
				}
				
				log.Fatalf("Failed to read incoming message. %s\n", err)
			}

			fmt.Printf("%s", msg)
		}
	}()

	// Read input and send to the server
	inputReader := bufio.NewReader(os.Stdin)
	for {
		msg, err := inputReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Fatalf("%s\n", msg)
			}

			fmt.Printf("Failed to read your message. %s\n", err)
			continue
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("Failed to write message to server. %s\n", err)
			continue
		}
	}
}