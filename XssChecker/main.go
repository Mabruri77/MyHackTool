package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	// URL target untuk diserang
	targetURL := "http://testphp.vulnweb.com/listproducts.php?cat=*"

	// Buka file payload
	payloadFile, err := os.Open("payload.txt")
	if err != nil {
		fmt.Println("Error opening payload file:", err)
		return
	}
	defer payloadFile.Close()

	// Buat scanner untuk membaca payload dari file
	scanner := bufio.NewScanner(payloadFile)

	// WaitGroup untuk menunggu selesainya semua goroutine
	var wg sync.WaitGroup

	// Mulai membaca payload dari file
	for scanner.Scan() {
		payload := scanner.Text()

		// Menyiapkan URL dengan parameter payload
		attackURL := strings.Replace(targetURL, "*", url.QueryEscape(payload), 1)

		// Menambahkan goroutine ke WaitGroup
		wg.Add(1)

		// Jalankan goroutine untuk mengirim request
		go func(url string, payload string) {
			defer wg.Done()
			result := checkPayload(url, payload)
			fmt.Println(result)
		}(attackURL, payload)
	}

	// Tunggu sampai semua goroutine selesai
	wg.Wait()
}

func checkPayload(url string, payload string) string {
	// Kirim request GET dengan payload
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Error sending request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Baca respons dari server
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response from %s: %v", url, err)
	}

	// Periksa apakah payload tercermin di respons
	if strings.Contains(string(body), payload) {
		return fmt.Sprintf("%sXSS Detected for payload: %s%s", ColorGreen+TextBold, payload, ColorReset)
	}
	return fmt.Sprintf("%sNo Vulnerability Detected for payload: %s%s", ColorRed, payload, ColorReset)
}
