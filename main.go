package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type User struct {
	Username string
	Conn net.Conn
	Send chan string
}

type Room struct {
	Users map[*User]bool
	Register chan *User
	Unregister chan *User
	Broadcast chan string
	mu sync.Mutex
}

func NewRoom() *Room {
	return &Room{
		Users: make(map[*User]bool),
		Register: make(chan *User),
		Unregister: make(chan *User),
		Broadcast: make(chan string),
	}
}

func (r *Room) Run() {
	for {
		select {
		case user := <-r.Register:
			r.mu.Lock()
			r.Users[user] = true
			fmt.Printf("%s has joined the chat.\r\n", user.Username)
			r.mu.Unlock()

		case user := <-r.Unregister:
			r.mu.Lock()
			if _, ok := r.Users[user]; ok {
				delete(r.Users, user)
				close(user.Send)
				fmt.Printf("%s has left the chat.\r\n", user.Username)
			}
			r.mu.Unlock()

		case msg := <-r.Broadcast:
			r.mu.Lock()
			for user := range r.Users {
				select {
				case user.Send<- msg:
				default:
					delete(r.Users, user)
					close(user.Send)
					fmt.Printf("Broadcast Error. %s has left the chat.\n", user.Username)
					fmt.Fprintf(user.Conn, "%s has left the chat.\r\n", user.Username)
				}
			}
			r.mu.Unlock()
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen to connection. %s\n", err)
	}
	defer listener.Close()

	room := NewRoom()
	go room.Run()

	fmt.Printf("Server running at localhost:8080\n")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection. %s\n" ,err)
			continue
		}

		go handleConnection(conn, room)
	}
}

func handleConnection(conn net.Conn, room *Room) {
	// Read username
	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read username. %s\r\n", err)
		return
	}
	username = strings.TrimSpace(username)

	// Register new user
	newUser := &User{
		Username: username,
		Conn: conn,
		Send: make(chan string),
	}
	room.Register<- newUser

	// Wait for everything to finish and disconnect user from room
	defer func () {
		room.Broadcast<- fmt.Sprintf("--- %s has left the chat.\r\n", newUser.Username)
		room.Unregister<- newUser
		conn.Write([]byte("You left the chat."))
		conn.Close()
	}()

	// Wait for incoming messages and display them to the user
	go func () {
		for msg := range newUser.Send {
			fmt.Fprintf(newUser.Conn, "%s", msg)
		}

		conn.Close()
	}()
	room.Broadcast<- fmt.Sprintf("--- %s has joined the chat.\r\n", newUser.Username)

	// Take input from user and broadcast to the room
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		room.Broadcast<- fmt.Sprintf("%s: %s\r\n", newUser.Username, msg)
	}
}