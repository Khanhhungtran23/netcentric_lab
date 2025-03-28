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
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server!")

	// create go routine to read response from the server 
	go func() { // is a paralel function
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() { // loop for each line in data from server sent to us
			fmt.Printf("Server: %s\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Lost connection to server: %v", err)
			os.Exit(1)
		}
	}()

	// Read input from user and send to server 
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages to send (type 'quit' to exit): ")
	// Loop for sending messages to Server
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		message := scanner.Text()
		if strings.TrimSpace(message) == "quit" {
			fmt.Println("Disconnecting...")
			break
		}

		// Sending message to server
		if _, err := conn.Write([]byte(message + "\n")); err != nil {
			log.Printf("Failed to send message to server")
			break 
		}
	}


}