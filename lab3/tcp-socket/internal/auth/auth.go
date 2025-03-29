package auth

import (
	"socket-tcp/internal/model"
	"encoding/base64"
	"sync"
	"errors"
	"crypto/rand"
	"math/big"
)

// AuthManager handles authentication and session management
// Using maps to quickly access following key	
type AuthManager struct {
	users 				map[string]*model.User				// username -> User
	connectedUsers		map[int]*model.ConnectedClient		// SessionID -> ConnectedClient
	mu 					sync.RWMutex 				// avoid race condition when many process access one resources - can be a variable
}

func NewAuthManager(users []*model.User) *AuthManager { // users slice
	// transform to format - map
	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.Username] = user 
	}

	return &AuthManager{
		users: 			userMap,
		connectedUsers: make(map[int]*model.ConnectedClient),
		mu:				sync.RWMutex{},
	}
}

// EncryptPassword - using base64
func EncryptPassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func VerifyPassword(plaintext, encrypted string) bool {
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return false
	}
	return string(decoded) == plaintext
}

func (am *AuthManager) AuthenticateUser(username, password string) (int, error) {
	am.mu.RLock()
	user, exists := am.users[username]
	am.mu.RUnlock()

	if !exists {
		return 0, errors.New("User not found")
	}

	if !VerifyPassword(password, user.Password) {
		return 0, errors.New("Invalid Password")
	}

	sessionID, err := GenerateSessionID()
	if err != nil {
		return 0, err
	}

	// Make sure we have unique sessionID
	am.mu.Lock()
	defer am.mu.Unlock()

	for {
		if _, exists := am.connectedUsers[sessionID]; !exists {
			break 
		}
		sessionID, err = GenerateSessionID()
		if err != nil {
			return 0, err
		}
	}
	// create and store connected client
	client := &model.ConnectedClient{
		User: user,
		SessionID: sessionID,
	}
	am.connectedUsers[sessionID] = client
	
	return sessionID, nil
}

// ValidateSession if session ID is valid
func (am *AuthManager) ValidateSession(sessionID int) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	_, exists := am.connectedUsers[sessionID]
	return exists
}

func GenerateSessionID() (int, error) {
	// Generate a random number between 100 & 999
	nBig, err := rand.Int(rand.Reader, big.NewInt(900))
	if err != nil {
		return 0, err
	}

	return int(nBig.Int64()) + 100, nil
}	