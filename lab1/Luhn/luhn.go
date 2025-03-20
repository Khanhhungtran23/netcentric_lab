package main

import (
	"fmt"
	"unicode"
)

func isValidLuhn(number string) bool {
	var digits []int

	// loc ra cac chi so tu chi so dau vao
	for _, char := range number {
		if unicode.IsDigit(char) {
			digits = append(digits, int(char-'0'))
		}
	}

	// neu so luong chu so nho hon 10, khong hop le
	if len(digits) <= 10 {
		fmt.Println("Credit number must be larger than 10")
		return false
	}

	sum := 0
	double := false

	for i := len(digits) - 1; i >= 0; i-- {
		n := digits[i]
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}

	return sum%10 == 0
}

func main() {
	testNumbers := []string{
		"4539 3195 0343 6467",
		"8273 1232 7352 0569",
		"79927398713",
	}

	for i := 0; i < len(testNumbers); i++ {
		if isValidLuhn(testNumbers[i]) {
			fmt.Println(testNumbers[i], "- This number is valid!")
		} else {
			fmt.Println(testNumbers[i], "- This number is not valid!")
		}
	}

	// for _, num := range testNumbers {
	// 	fmt.Print("Số %s hợp lệ? %s\n", num, isValidLuhn(num))
	// }
}
