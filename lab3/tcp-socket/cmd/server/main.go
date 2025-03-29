package main

import (
	"fmt" // lib for function to print out to the screen - Println - Printf
	"log" // for logging
	"net"  // tcp lib - dependency : For functions working with network like TCP, UDP, HTTP, ...
	"strings"
	"flag"
	"os"

	"socket-tcp/internal/protocol"
	"socket-tcp/internal/auth"
	"socket-tcp/internal/storage"
)


var (
	port 		= flag.String("port", "8080", "Server port")
	userFile	= flag.String("users", "data/users.json", "User data file")
	storageType	= flag.String("storage", "json", "Storage type (json or gob)")
)



// main func to run	
func main() {
	// Parse cmd-line flags
	flag.Parse()

	// Set up logger
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)


	// Initialize the storage
	var st storage.StorageType 
	if *storageType == "gob" {
		st = storage.GOBStorage
	} else {
		st = storage.JSONStorage
	}

	// create user storage
	userStorage := storage.NewUserStorage(*userFile, st)
	// Load Users
	users, err := userStorage.LoadUsers()
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	log.Printf("Loaded %d users from file", len(users))

	// Create auth manager
	authManager := auth.NewAuthManager(users)

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to start tcp server: %v" , err)
	}
	// delay to execute this command until func main ending
	// help the resources free if the errors occur
	defer listener.Close()

	fmt.Printf("Server TCP is running on port %s!\n", *port)

	// Accept and Handle Connecting
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accepting connect from : %v", err)
			continue
		}

		// handle connect in each goroutine
		go handleConnection(conn, authManager)
	}
}

func handleConnection(conn net.Conn, authManager *auth.AuthManager) {
	defer conn.Close() // close connect when this function ending to avoid resource leakage

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in connection handler: %v", r)
		}
	}()


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
		fmt.Printf("Received from %s: Command= %s - SessionID= %d - Payload= %s\n", clientAddr, msg.Command, msg.SessionID, msg.Payload)
		
		// Process message based on commamd
		switch msg.Command {
		case protocol.CmdAuth:
			// Handle authentication
			if authenticated {
				if err := msgHandler.SendMessage(sessionID, protocol.CommandType("ERROR"), "Already authenticated"); err != nil {
					log.Printf("Failed to send error message: %v", err)
				}
				continue // ignore new cmd line
			}

			// Virtual authenticate simply
			parts := strings.SplitN(msg.Payload, " ", 2)
			if len(parts) != 2 {
				msgHandler.SendMessage(0, protocol.CommandType("ERROR"), "Invalid auth format")
				continue
			}

			username, password := parts[0], parts[1]

			// Authenticate user
			newSessionID, err := authManager.AuthenticateUser(username, password)
			if err != nil {
				if err := msgHandler.SendMessage(0, protocol.CommandType("ERROR"), "Authentication Failed: " + err.Error()); err != nil {
					log.Print("Failed to send error message: %v", err)
				}
				continue
			}

			sessionID = newSessionID
			authenticated = true

			if err := msgHandler.SendMessage(sessionID, protocol.CommandType("OK"), fmt.Sprintf("Authentication Successful. Your session ID is %d", sessionID)); err != nil {
				log.Printf("Failed to send success message: %v", err)
				break
			}
			log.Printf("Client %s authenticated as %s with session ID %d", clientAddr, username, sessionID)

		case protocol.CmdQuit:
			if authenticated && msg.SessionID != sessionID{
				if err := msgHandler.SendMessage(sessionID, protocol.CommandType("ERROR"), "Invalid session ID"); err != nil {
					log.Printf("Failed to send error message: %v", err)
				}
				continue
			}

			// Send goodbye
			if authenticated {
				if err := msgHandler.SendMessage(sessionID, protocol.CommandType("BYE"), "Goodbye!"); err != nil {
					log.Printf("Failed to send goodbye message: %v", err)
				}
			} else {
				if err := msgHandler.SendMessage(0, protocol.CommandType("BYE"), "Goodbye!"); err != nil {
					log.Printf("Failed to send goodbye message: %v", err)
				}
			}
			
			log.Printf("Client %s disconnected", clientAddr)
			return
		default:
			// Check authentication
			if !authenticated {
				if err := msgHandler.SendMessage(0, protocol.CommandType("ERROR"), "Not authenticated"); err != nil {
					log.Printf("Failed to send error message: %v", err)
				}
				continue
			}

			// Check session
			if msg.SessionID != sessionID {
				if err := msgHandler.SendMessage(sessionID, protocol.CommandType("ERROR"), "Invalid session ID"); err != nil {
					log.Printf("Failed to send error message: %v", err)
				}
				continue
			}

			// Handle other commands (will implement later)
			if err := msgHandler.SendMessage(sessionID, protocol.CommandType("ECHO"), 
				fmt.Sprintf("Received command: %s with payload: %s", msg.Command, msg.Payload)); err != nil {
				log.Printf("Failed to send response message: %v", err)
				break
			}
		}
	}

	log.Printf("Connection from %s closed!", clientAddr)
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