package main

import (
	"fmt" // lib for function to print out to the screen - Println - Printf
	"log" // for logging
	"net"  // tcp lib - dependency : For functions working with network like TCP, UDP, HTTP, ...
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

	conn.Write([]byte("Welcome to TCP server!\n")) // send welcome message to client - transform into bytes because Write method require data as byte format

	buffer := make([]byte, 1024) // create a slice as byte with size 1024 to contain data from client

	for {
		n, err := conn.Read(buffer) // read data from client into Buffer, return reading number of bytes n
		if err != nil {
			fmt.Printf("Connection from %s closed: %v \n", clientAddr, err)
			return 
		}

		// showing data getting from client
		data := buffer[:n]
		fmt.Printf("Recieved from %s: %s", clientAddr, string(data)) // byte to string

		// Sending repsonse Echo
		conn.Write([]byte(fmt.Sprintf("Echo: %s", string(data))))

	}
	
}