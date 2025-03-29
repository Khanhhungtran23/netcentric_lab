package model

// User struct - represent user in the system
type User struct {
	Username string 		`json:"username"`
	Password string			`json:"password"` // Encrypted
	Fullname string			`json:"fullname"`
	Emails    []string 		`json:"emails"`
	Addresses  []Address	`json:"addresses"`
}

type Address struct {
	Type		string		`json:"type"`
	Details		string 		`json:"details"`
}

type ConnectedClient struct {
	User 					*User
	SessionID 				int	// unique random key
}

// type GameState struct {
// 	Target 			int 
// 	GuessCount 		int
// 	InProgress		bool 
// }

