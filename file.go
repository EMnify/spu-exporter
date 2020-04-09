package main

import (
	"bufio"
	"log"
	"os"
)

// Deprecated, just for testing without using ssh. Might be useful for creating tests.

func readFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
