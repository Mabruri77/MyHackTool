package code

import (
	"bufio"
	"os"
)

func readPayloadsFromFile(filename string) ([]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var payloads []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		payload := scanner.Text()
		if payload != "" {
			payloads = append(payloads, payload)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return payloads, nil
}
