/* Protocol - This is the convention about format and decoding message between client & server
Although they already connect & communicate with TCP but need a clear structue format to :
	1. Determine type of command (AUTH, GUESS, FILE, v.v)
	2. Attach SessionID (exercise requirements)
	3. Extract the payload 
In this scene, this protocol help both side can understand the message in a same way. It acts like a shared language for both client and server be able to understand in the right way.
Instead of text communicate without any structure.
*/  

package protocol

import (
	"fmt"
	"net"
	"strings"
	"errors"
	"strconv"
	"bufio"
)

type CommandType string // Define type of command - just a field/attribute in string 

const (
	CmdAuth 		CommandType = "AUTH"
	CmdFile 		CommandType = "FILE"
	CmdGuess 		CommandType = "GUESS"
	CmdQuit 		CommandType = "QUIT" // quit the connection
	CmdStartGame 	CommandType = "START"
	CmdEndGame 		CommandType = "END"
)

// Define format of message - a wrapper
type Message struct {
	SessionID 		int 
	Command 		CommandType
	Payload 		string // of data wanna send
}

// Method to send and receive message
type MessageHandler struct {
	conn net.Conn
}

// Create a new MessageHandler to new MessageHandler
func NewMessageHandler(conn net.Conn) *MessageHandler {
	return &MessageHandler {
		conn: conn,
	}
}

// This is a func having receiver in Go, specific is struct MessageHandler
// It is use to send message through TCP connection
// (mh *MessageHandler is the receiver of func) => It means func SendMessage is belongs to MessageHandler
// mh is the presentation variable for object MessageHandler
// *MessageHandler is the pointer helping func can be able to change data in struct if needed.
func (mh *MessageHandler) SendMessage(sessionID int, command CommandType, payload string) error {
	message := ""

	if command == CmdAuth {
		message = fmt.Sprintf("%s %s\n", command, payload)
	} else {
		message = fmt.Sprintf("%d_%s %s\n", sessionID, command, payload)
	}

	_, err := mh.conn.Write([]byte(message))
	return err
}

func (mh *MessageHandler) ReadMessage() (*Message, error) {
	// read message from server
	reader := bufio.NewReader(mh.conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)

	var sessionID int
	var commandStr string
	var payload string

	// check if having the cmdAuth or not ?
	if strings.HasPrefix(line, string(CmdAuth)) {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			return nil, errors.New("Format message is not valid")
		}
		commandStr = parts[0]
		payload = parts[1]
	} else {
		// Separate session ID and cmd
		parts := strings.SplitN(line, "_", 2)
		if len(parts) < 2 {
			return nil, errors.New("Format message is not valid")
		}
		var err error
		sessionID, err = strconv.Atoi(parts[0]) // convert string into integer - Atoi
		if err != nil {
			return nil, errors.New("Session ID is not valid")
		}

		cmdParts := strings.SplitN(parts[1], " ", 2)
		if len(cmdParts) < 2 {
			return nil, errors.New("Format message is not valid")
		}
		commandStr = cmdParts[0]
		payload = cmdParts[1]
	}

	// init object Message
	message := &Message{
		SessionID: sessionID,
		Command:   CommandType(commandStr),
		Payload:   payload,
	}

	return message, nil
}

