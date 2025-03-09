package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"time"
)

var websites = []string{
	"https://google.com",
	"https://facebook.com",
	"https://1212121212111212122.org",
}

// app password need to be generated from https://myaccount.google.com/apppasswords
const (
	smtpServer = "smtp.gmail.com"
	smtpPort   = "587"
	sender     = "murtadhabg@gmail.com"
	password   = "udgbujilpmtehtpf "
	receiver   = "murtadhabg@gmail.com"
)

var downSince = make(map[string]time.Time)

func sendEmail(site string, downtimeDuration time.Duration) {
	subject := "Website Down Alert"
	body := fmt.Sprintf("The website %s is down. %s\n Down since: %v", site, time.Now().Format(time.RFC1123), downtimeDuration)
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", sender, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, sender, []string{receiver}, []byte(message))

	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("📧 Alert email sent")
	}
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
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	var logMessage string

	if err != nil {
		if _, exist := downSince[url]; !exist {
			downSince[url] = time.Now()
		}
		downtimeDuration := time.Since(downSince[url])
		logMessage = fmt.Sprintf("%s ❌  Website is down [%s] - Was down for: %v\n", url, time.Now().Format(time.RFC1123), downtimeDuration)
		sendEmail(url, downtimeDuration)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			if _, exist := downSince[url]; exist {
				delete(downSince, url)
			}
			logMessage = fmt.Sprintf("%s ✅  Website is up [%s] - Response time: %v\n", url, time.Now().Format(time.RFC1123), elapsed)
		} else {
			logMessage = fmt.Sprintf("%s ⚠️  website status is: %d [%s] - Response time: %v\n", url, resp.StatusCode, time.Now().Format(time.RFC1123), elapsed)
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
