# TCP-based CLI Chat Application

### Description
A simple TCP-based chat application and client built in Go, able to handle multiple users simultaneously in a same connection room.

### Purpose
This project was built with the goal of understanding networking concepts, with a practical application of how TCP works. This project also helped me understand goroutines and channels strong capabilities to help handle connections and communication.

### Installation/How to run
- Run the server:
```
go run main.go
```

- Run the client:
First, navigate to the client directory.
```
go run client.go
```

### Usage
By running the client, the user will directly connect to the server, and be prompted to add a username. The server will notify other users that a new user has connected, and all connected users will be able to communicate. By closing the terminal, the user will disconnect from the server, and all still connected users will be notified.

### Future Improvements and Changes
Right now, this is a simple study project to understand networking concepts and TCP. I want to improve this project and allow users to create their own rooms with connection codes, where any user with the code can connect.

### License
MIT License