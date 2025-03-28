package model

// User struct - represent user in the system
type User Struct {
	Username string 		`json:"username"`
	Password string			`json:"password"` // Encrypted
	Fullname string			`json:"fullname"`
	Emails    []string 		`json:"emails"`
	Addresses  []Address	`json:"addresses"`
}

type Address Struct {
	Type		string		`json:"type"`
	Details		string 		`json:"details"`
}

type ConnectedClient Struct {
	User 					*User
	SessionID 				int	// unique random key
}

type GameState Struct {
	Target 			int 
	GuessCount 		int
	InProgress		bool 
}

