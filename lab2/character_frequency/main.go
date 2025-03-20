package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime" // access runtime info - CPU cores
	"sort"
	"sync" // sync goroutine
)

// CharFrequency represents a character and its frequency
type CharFrequency struct {
	Char  rune // unicode representation
	Count int
}

// countCharacters counts character frequencies in a text chunk
func countCharacters(text string) map[rune]int {
	freqMap := make(map[rune]int)
	for _, char := range text {
		freqMap[char]++
	}
	return freqMap
}

// processTextConcurrently splits text into chunks and processes them concurrently
func processTextConcurrently(text string) map[rune]int {
	// Determine number of goroutines based on CPU cores
	numCPU := runtime.NumCPU()
	chunkSize := (len(text) + numCPU - 1) / numCPU // Round up division

	// Create channels for communication
	results := make(chan map[rune]int, numCPU)
	var wg sync.WaitGroup

	// Process text in chunks
	for i := 0; i < len(text); i += chunkSize {
		wg.Add(1)
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}

		chunk := text[i:end]
		go func(chunk string) {
			defer wg.Done()
			chunkResult := countCharacters(chunk)
			results <- chunkResult
		}(chunk)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Combine results
	finalResult := make(map[rune]int)
	for result := range results {
		for char, count := range result {
			finalResult[char] += count
		}
	}

	return finalResult
}

// readFile reads text from a file
func readFile(filePath string) (string, error) {
	// Resolve path - if relative, convert to absolute
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("Wrong path of file: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("File does not exist: %s", absPath)
	}

	// Open file
	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("Cannot open file: %v", err)
	}
	defer file.Close()

	// Read file content
	scanner := bufio.NewScanner(file)
	var text string
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error when reading file: %v", err)
	}

	return text, nil
}

// printFrequencies prints character frequencies in sorted order
func printFrequencies(freqMap map[rune]int) {
	// Convert map to slice for sorting
	var freqList []CharFrequency
	for char, count := range freqMap {
		freqList = append(freqList, CharFrequency{char, count})
	}

	// Sort by character
	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Char < freqList[j].Char
	})

	// Print frequencies
	for _, cf := range freqList {
		if cf.Char == ' ' {
			fmt.Printf("(blank): %d\n", cf.Count)
		} else if cf.Char == '\n' {
			fmt.Printf("\\n: %d\n", cf.Count)
		} else {
			fmt.Printf("%c: %d\n", cf.Char, cf.Count)
		}
	}
}

func main() {
	var text string
	var err error

	// Check if file path is provided
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		text, err = readFile(filePath)
		if err != nil {
			fmt.Printf("Errors: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Reading file: %s\n", filePath)
	} else {
		fmt.Println("No file provided. Please enter a string:")
		reader := bufio.NewReader(os.Stdin)
		text, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		text = text[:len(text)-1] // Loại bỏ ký tự xuống dòng '\n'
	}

	// Process text concurrently
	frequencies := processTextConcurrently(text)

	// Print results
	fmt.Println("Result:")
	printFrequencies(frequencies)
}
