package main

import (
	"bufio"
	"fmt"
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

	// Read input and send to server
	go func () {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil { 
				/* 
				TO FIX: WHEN CLOSING THE SERVER WITHOUT ENDING CLIENT SIDE
				- THE CURRENT FORMAT CAUSES THE ERROR MESSAGE TO PRINT IN A LOOP
				- RETURN WILL CAUSE THE MESSAGE TO PLAY ONCE, BUT CLIENT WILL STILL RUN
				- LOG.FATALF SOLVES THE PROBLEM, BUT TRY TO IMPLEMENT A WAY WHERE LOG.FATALF WILL ONLY BE USED FOR WHEN THE SERVER AND CLIENT CONNECTION IS TERMINATED
				*/
				fmt.Printf("Failed to read incoming message. %s\n", err)
				continue
			}

			fmt.Printf("%s", msg)
		}
	}()

	// Receive input from server and print
	inputReader := bufio.NewReader(os.Stdin)
	for {
		msg, err := inputReader.ReadString('\n')
		if err != nil {
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