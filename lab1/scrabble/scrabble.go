package main

import (
	"fmt"
	"strings"
)

// func to calculate Scrabble score for word or sentences
func ScrabbleScore(word string) int {
	scrabbleValues := map[rune]int{
		'A': 1, 'E': 1, 'I': 1, 'O': 1, 'U': 1, 'L': 1, 'N': 1, 'R': 1, 'S': 1, 'T': 1,
		'D': 2, 'G': 2,
		'B': 3, 'C': 3, 'M': 3, 'P': 3,
		'F': 4, 'H': 4, 'V': 4, 'W': 4, 'Y': 4,
		'K': 5,
		'J': 8, 'X': 8,
		'Q': 10, 'Z': 10,
	}
	totalScore := 0

	word = strings.ToUpper(word)

	for _, letter := range word {
		if val, exists := scrabbleValues[letter]; exists {
			totalScore += val
		}
	}
	return totalScore
}

func main() {
	var input string
	fmt.Println("Input your word or string for Scrabble Score:")
	fmt.Scanln(&input)

	score := ScrabbleScore(input)
	fmt.Printf("Scrabble Score của '%s' là: %d\n", input, score)
}
