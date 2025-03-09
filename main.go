package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"time"
)

type Config struct {
	Websites      []string      `yaml:"websites`
	CheckInterval time.Duration `yaml:"check_interval"`
	Email         struct {
		SMTPServer string `yaml:"smtpServer"`
		SMTPPort   string `yaml:"smtpPort"`
		Sender     string `yaml:"sender"`
		Password   string `yaml:"password"`
		Receiver   string `yaml:"receiver"`
	} `yaml:"email"`
}

var config Config

func loadConfig() error {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, &config)
}

var downSince = make(map[string]time.Time)

func sendEmail(site string, downtimeDuration time.Duration) {
	subject := "Website Down Alert"
	body := fmt.Sprintf("The website %s is down. %s\n Down since: %v", site, time.Now().Format(time.RFC1123), downtimeDuration)
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", config.Email.Sender, config.Email.Password, config.Email.SMTPServer)
	err := smtp.SendMail(config.Email.SMTPServer+":"+config.Email.SMTPPort, auth, config.Email.Sender, []string{config.Email.Receiver}, []byte(message))

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
	err := loadConfig()

	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	for {
		fmt.Println("\n --- Checking websites...")
		var wg sync.WaitGroup
		for _, url := range config.Websites {
			wg.Add(1)
			go checkWebsite(url, &wg)
		}
		wg.Wait()
		fmt.Println("--- Waiting for the next check ... ---")
		time.Sleep(config.CheckInterval)
	}
}
