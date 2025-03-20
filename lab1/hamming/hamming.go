package main

import (
	"fmt"
	"math/rand"
	"time"
)

// calculate hamming distance between 2 dna strands with same length
func HammingDistance(dna1, dna2 string) int {
	if len(dna1) != len(dna2) {
		panic("DNA sequences must have the same lennth!")
	}

	distance := 0
	for i := 0; i < len(dna1); i++ {
		if dna1[i] != dna2[i] {
			distance++
		}
	}
	return distance
}

// generate sequences DNA randomly with specific n length
func RandomDNA(length int) string {
	dnaBases := "ACGT"
	rand.Seed(time.Now().UnixNano())

	dna := make([]byte, length)
	for i := 0; i < length; i++ {
		dna[i] = dnaBases[rand.Intn(len(dnaBases))]
	}
	return string(dna)
}

func main() {
	dna1 := "GAGCCTACTAACGGGAT"
	dna2 := "CATCGTAATGACGGCCT"
	fmt.Println("Test with 2 common sample: ")
	fmt.Println("DNA 1:", dna1)
	fmt.Println("DNA 2:", dna2)
	fmt.Println("Hamming Distance: ", HammingDistance(dna1, dna2))

	const numTests = 1000
	const dnaLength = 18
	fmt.Println("\nRunning 1000 random DNA tests...")

	for i := 0; i < numTests; i++ {
		randDNA1 := RandomDNA(dnaLength)
		randDNA2 := RandomDNA(dnaLength)
		result := HammingDistance(randDNA1, randDNA2)
		fmt.Println("Distance ", i, ": ", result)
	}
	fmt.Println("Completed 1000 random DNA tests.")
}
