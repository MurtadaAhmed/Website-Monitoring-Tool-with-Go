package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var websites = []string{
	"https://google.com",
	"https://facebook.com",
	"https://1212121212111212122.org",
}

func checkWebsite(url string, wg *sync.WaitGroup) {

	defer wg.Done()

	// 1. reading the file to log the output
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 2. checking the url
	resp, err := http.Get(url)

	var logMessage string

	if err != nil {
		logMessage = fmt.Sprintf("%s ❌  Website is down [%s]\n", url, time.Now().Format(time.RFC1123))
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			logMessage = fmt.Sprintf("%s ✅  Website is up [%s]\n", url, time.Now().Format(time.RFC1123))
		} else {
			logMessage = fmt.Sprintf("%s ⚠️  website status is: %d [%s]", url, resp.StatusCode, time.Now().Format(time.RFC1123))
		}
	}

	fmt.Println(logMessage)

	// 3. saving the output to the log file
	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func main() {
	interval := 30 * time.Second

	for {
		fmt.Println("\n --- Checking websites...")
		var wg sync.WaitGroup
		for _, url := range websites {
			wg.Add(1)
			go checkWebsite(url, &wg)
		}
		wg.Wait()
		fmt.Println("--- Waiting for the next check ... ---")
		time.Sleep(interval)
	}
}
