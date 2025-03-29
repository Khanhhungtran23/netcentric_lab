package main

import (
	"fmt" // lib for function to print out to the screen - Println - Printf
	"log" // for logging
	"net"  // tcp lib - dependency : For functions working with network like TCP, UDP, HTTP, ...
	"strings"

	"socket-tcp/internal/protocol"
)

// main func to run	
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start tcp server: %v" , err)
	}
	// delay to execute this command until func main ending
	// help the resources free if the errors occur
	defer listener.Close()

	fmt.Println("Server TCP is running on port 8080!")

	// Accept and Handle Connecting
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accepting connect from : %v", err)
			continue
		}

		// handle connect in each goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // close connect when this function ending to avoid resource leakage
	clientAddr := conn.RemoteAddr().String() // get IP address and Port of client
	fmt.Printf("New connection from %s\n", clientAddr) // on terminal

	// conn.Write([]byte("Welcome to TCP server!\n")) // send welcome message to client - transform into bytes because Write method require data as byte format
	msgHandler := protocol.NewMessageHandler(conn)
	if err := msgHandler.SendMessage(0, protocol.CommandType("SERVER"), "Welcome to TCP Socket Server! Please use AUTH username password to login."); err != nil {
			log.Printf("Failed to send welcome message: %v", err)
			return 
	}
	
	// Default sessionID
	sessionID := 0
	authenticated := false

	// [process to handle receive message from client]
	for {
		msg, err := msgHandler.ReadMessage()
		if err != nil {
			fmt.Printf("Connection from %s closed: %v \n", clientAddr, err)
			return
		}

		// Showing message received
		fmt.Printf("Received from %s: Command=%s - SessionID=%d - Payload=%s\n", clientAddr, msg.Command, msg.SessionID, msg.Payload)
		switch msg.Command {
		case protocol.CmdAuth:
			// Virtual authenticate simply
			parts := strings.SplitN(msg.Payload, " ", 2)
			if len(parts) != 2 {
				msgHandler.SendMessage(0, protocol.CommandType("ERROR"), "Invalid auth format")
				continue
			}

			username, password := parts[0], parts[1]

			if username == "admin" && password == "123" {
				sessionID = 333
				authenticated = true
				msgHandler.SendMessage(sessionID, protocol.CommandType("OK"), fmt.Sprintf("Authentication Successful. Your sessionID is %d", sessionID))
			} else {
				msgHandler.SendMessage(0, protocol.CommandType("ERROR"), "Authentication Fail")
			}
		case protocol.CmdQuit:
			if authenticated {
				msgHandler.SendMessage(sessionID, protocol.CommandType("BYE"), "See Ya!")
			} else {
				msgHandler.SendMessage(0, protocol.CommandType("BYE"), "You are not login to quit")
			}
			return
		default:
			if !authenticated {
				msgHandler.SendMessage(0, protocol.CommandType("ERROR"), "Not authenticated")
				continue
			}

			if msg.SessionID != sessionID {
				msgHandler.SendMessage(sessionID, protocol.CommandType("ERROR"), "Invalid session ID")
				continue
			}

			// Send default response echo
			msgHandler.SendMessage(sessionID, protocol.CommandType("ECHO"), fmt.Sprintf("Echo: %s", msg.Payload))
		}
	}






	// buffer := make([]byte, 1024) // create a slice as byte with size 1024 to contain data from client

	// for {
	// 	n, err := conn.Read(buffer) // read data from client into Buffer, return reading number of bytes n
	// 	if err != nil {
	// 		fmt.Printf("Connection from %s closed: %v \n", clientAddr, err)
	// 		return 
	// 	}

	// 	// showing data getting from client
	// 	data := buffer[:n]
	// 	fmt.Printf("Recieved from %s: %s", clientAddr, string(data)) // byte to string

	// 	// Sending repsonse Echo
	// 	conn.Write([]byte(fmt.Sprintf("Echo: %s", string(data))))

	// }
	
}