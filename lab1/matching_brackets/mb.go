package main

import (
	"fmt"
)

func isValidBrackets(s string) bool {
	stack := []rune{}
	bracketMap := map[rune]rune{')': '(', '}': '{', ']': '['}

	for _, char := range s {
		switch char {
		case '(', '{', '[':
			stack = append(stack, char)
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != bracketMap[char] {
				return false
			}
			// pop from stack
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func main() {
	// Các test case
	tests := []string{
		"([]{})",
		"([)]",
		"{[()()]}",
		"{{[[(())]]}}",
		"{[}",
	}

	for _, test := range tests {
		fmt.Printf("String \"%s\" -> %v\n", test, isValidBrackets(test))
	}
}
