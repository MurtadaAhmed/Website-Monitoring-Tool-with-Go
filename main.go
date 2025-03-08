package main

import (
	"fmt"
	"net/http"
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

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(url, " ❌  Website is down", time.Now())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println(url, "✅  website is up", time.Now())
	} else {
		fmt.Println(url, "⚠️  website status is ", resp.StatusCode, time.Now())
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
