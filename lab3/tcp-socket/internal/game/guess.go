package game 

import (
	"sync"

	"socket-tcp/internal/model"
)

type GuessingGame struct {
	games 		map[int]*model.GameState	// sessionID -> GameState
	mu 			sync.RWMutex
}

func NewGuessingGame() *GuessingGame {
	return &GuessingGame{
		games: 	make(map[int]*model.GameState),
		mu: 	sync.RWMutex{},
	}
}

func (gg *GuessingGame) StartGame(sessionID int) (string, error) {

}

func (gg *GuessingGame) MakeGuess() {

}

func (gg *GuessingGame) EndGame() {

}

// has active game trakc if a session has an active game
func (gg *GuessingGame) HasActiveGame(sessionID int) bool {

}