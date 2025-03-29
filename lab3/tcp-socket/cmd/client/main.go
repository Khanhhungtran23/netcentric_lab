package main

import (
	"fmt"
	"net"
	"log"
	"bufio" // read input data from keyboard and from network connecting efficiently
	"os" // interact with the system like exit program
	"strings" // handle string (trim white space - check confition)
	// trim : to remove white space at start and end of the string (can be at a specific-word)
	// In Golang: from strings - providing Trim() and TrimSpace()
	// Also having TrimLeft - TrimRight

	"socket-tcp/internal/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server!")

	msgHandler := protocol.NewMessageHandler(conn)
	msgHandler.SendMessage(0, protocol.CommandType("GREET"), "Hello from Khanh Hung\n")
	// conn.Write([]byte("Hello Server from KhanhHung!\n"))

	var sessionID 		int
	var authenticated 	bool

	// create go routine to read response from the server 
	go func() { // is a paralel function

		for {
			msg, err := msgHandler.ReadMessage()
			if err != nil {
				fmt.Println("Lost connection to Server: %v", err)
				os.Exit(1)
				return 
			}

			// Process Message based on type
			switch msg.Command {
			case protocol.CommandType("OK"):  
				if strings.Contains(msg.Payload, "Authentication Succesful") {
					sessionID = msg.SessionID
					authenticated = true 
					fmt.Printf("\nAuthenticated with session ID: %d\n", sessionID)
				}
				fmt.Printf("\nServer: %s\n", msg.Payload)
			case protocol.CommandType("ERROR"), protocol.CommandType("SERVER"), protocol.CommandType("ECHO"):
				fmt.Printf("\nServer: %s\n", msg.Payload)
			
			case protocol.CommandType("BYE"):
				fmt.Printf("\nServer: %s\n", msg.Payload)
				os.Exit(0)
			default:
				fmt.Printf("\nServer [%s]: %s\n", msg.Command, msg.Payload)
			}
			// fmt.Print("> ")
		}
		// scanner := bufio.NewScanner(conn)
		// for scanner.Scan() { // loop for each line in data from server sent to us
		// 	fmt.Printf("Server: %s\n", scanner.Text())
		// }

		// if err := scanner.Err(); err != nil {
		// 	fmt.Println("Lost connection to server: %v", err)
		// 	os.Exit(1)
		// }
	}()

	// Read input from user and send to server 
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println("Type messages to send (type 'quit' to exit): ")
	// Loop for sending messages to Server
	// for {
	// 	fmt.Print("> ")
	// 	if !scanner.Scan() {
	// 		break
	// 	}

	// 	message := scanner.Text()
	// 	if strings.TrimSpace(message) == "quit" {
	// 		fmt.Println("Disconnecting...")
	// 		break
	// 	}

	// 	// Sending message to server
	// 	if _, err := conn.Write([]byte(message + "\n")); err != nil {
	// 		log.Printf("Failed to send message to server")
	// 		break 
	// 	}
	// }

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please type 'help' for available commands: ")
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() { break }

		input := scanner.Text()
		if input == "" {
			continue
		} 

		// Analyze commands
		parts := strings.SplitN(input, " ", 2)
		command := strings.ToUpper(parts[0])
		payload := ""

		if len(parts) > 1 {
			payload = parts[1]
		}

		switch command {
		case "HELP":
			fmt.Println("Available Commands: ")
			fmt.Println("  Auth username  password  - Authentication with the server")
			fmt.Println("  QUIT					    - Disconnect from the server")
			if authenticated {
				fmt.Println("  START 				- Start a new guessing game")
				fmt.Println("  GUESS number 		- Make a guess")
				fmt.Println("  END 					- End the current game")
				fmt.Println("  FILE filename		- Download a file")
			}

		case "AUTH":
			if authenticated {
				fmt.Println("Already Authenticated")
				continue 	
			}
			err := msgHandler.SendMessage(0, protocol.CmdAuth, payload)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
                return
			}
		case "QUIT":
            var err error
            if authenticated {
                err = msgHandler.SendMessage(sessionID, protocol.CmdQuit, "")
            } else {
                err = msgHandler.SendMessage(0, protocol.CmdQuit, "")
            }
            if err != nil {
                log.Printf("Failed to send message: %v", err)
            }
            return

        default:
            if !authenticated {
                fmt.Println("Not authenticated. Use AUTH username password")
                continue
			}
			
			var cmdType protocol.CommandType 
			switch command {
			case "START":
                cmdType = protocol.CmdStartGame
            case "GUESS":
                cmdType = protocol.CmdGuess
            case "END":
                cmdType = protocol.CmdEndGame
            case "FILE":
                cmdType = protocol.CmdFile
            default:
                fmt.Println("Unknown command. Type 'help' for available commands")
                continue
			}

			// Send command with Session ID
			err := msgHandler.SendMessage(sessionID, cmdType, payload)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
		}
	}
}