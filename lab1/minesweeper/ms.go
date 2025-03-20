package main

import (
	"fmt"
	"math/rand"
	"time"
)

// size of fields
const rows, cols = 20, 25
const mineCount = 99

// func to create random mines
func generateMinefield() [][]rune {
	board := make([][]rune, rows)
	for i := range board {
		board[i] = make([]rune, cols)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}
	rand.Seed(time.Now().UnixNano())
	for mines := 0; mines < mineCount; {
		r, c := rand.Intn(rows), rand.Intn(cols)
		if board[r][c] != '*' {
			board[r][c] = '*'
			mines++
		}
	}
	return board
}

// Hàm đếm số mìn xung quanh 1 ô
func countMines(board [][]rune, r, c int) int {
	count := 0
	directions := []struct{ dr, dc int }{
		{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1},
	}

	for _, d := range directions {
		nr, nc := r+d.dr, c+d.dc
		if nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] == '*' {
			count++
		}
	}
	return count
}

// Hàm cập nhật bảng với số lượng mìn xung quanh mỗi ô
func updateBoard(board [][]rune) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if board[r][c] == '.' { // Nếu là ô trống, tính số mìn xung quanh
				mineCount := countMines(board, r, c)
				if mineCount > 0 {
					board[r][c] = rune('0' + mineCount)
				}
			}
		}
	}
}

// Hàm in bảng
func printBoard(board [][]rune) {
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Println()
	}
}

func main() {
	// Tạo và cập nhật bảng mìn
	board := generateMinefield()
	updateBoard(board)

	// In bảng kết quả
	printBoard(board)
}
