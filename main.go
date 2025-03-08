package main

import (
	"fmt"
	"net/http"
	"time"
)

var websites = []string{
	"https://google.com",
	"https://facebook.com",
}

func checkWebsite(url string) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Website is down", time.Now())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("website is up", time.Now())
	} else {
		fmt.Println("website status is ", resp.StatusCode, time.Now())
	}
}

func main() {
	interval := 30 * time.Second

	for {
		fmt.Println("\n --- Checking website...")
		for index, url := range websites {
			fmt.Println("[", index+1, "] ", url)
			checkWebsite(url)
		}
		fmt.Println("--- Waiting for the next check ... ---")
		time.Sleep(interval)
	}
}
