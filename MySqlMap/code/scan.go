package code

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

func SqlMap() {
	// Target URL to test
	targetURL := "http://testphp.vulnweb.com/listproducts.php"

	// Filename containing payloads
	payloadsFile := "payloads.txt"

	// Read payloads from file
	payloads, err := readPayloadsFromFile(payloadsFile)
	if err != nil {
		fmt.Println("Error reading payloads from file:", err)
		return
	}

	// Wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Iterate over each payload
	for _, payload := range payloads {
		wg.Add(1) // Increment the wait group counter
		go func(p string) {
			defer wg.Done() // Decrement the wait group counter when the goroutine finishes
			// Send HTTP request with payload
			startTime := time.Now()
			resp, err := http.Get(targetURL + "?cat=" + p)
			if err != nil {
				fmt.Println("Error sending HTTP request:", err)
				return
			}
			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Non-OK status code (%d) for payload: %s\n", resp.StatusCode, p)
				return
			}

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			// Check for potential SQL injection indicators in response body
			if strings.Contains(string(body), "SQL syntax") || strings.Contains(string(body), "mysql_fetch_array") {
				fmt.Printf("Potential SQL injection vulnerability detected in response body for payload: %s\n", p)
			}

			// Check response time for potential blind SQL injection (time-based)
			elapsedTime := time.Since(startTime)
			if strings.Contains(p, "SLEEP") && elapsedTime > time.Second*4 {
				fmt.Printf("Potential blind SQL injection vulnerability (time-based) detected with payload: %s\n", p)
			}

			// Check response headers for potential anomalies
			for header, values := range resp.Header {
				for _, value := range values {
					if strings.Contains(value, "error") || strings.Contains(value, "mysql") {
						fmt.Printf("Potential SQL injection vulnerability detected in response header (%s: %s) for payload: %s\n", header, value, p)
					}
				}
			}

			// Check response body for potential anomalies
			if strings.Contains(string(body), "error") || strings.Contains(string(body), "MySQL") {
				fmt.Printf("Potential SQL injection vulnerability detected in response body for payload: %s\n", p)
			}

			// Add more advanced checks as needed...
		}(payload) // Pass the payload to the anonymous function
	}
	fmt.Println("please wait....")
	wg.Wait() // Wait for all goroutines to finish before returning
}
