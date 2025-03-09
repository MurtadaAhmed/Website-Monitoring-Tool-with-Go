package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "monitor.db")
	if err != nil {
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS logs (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	website TEXT,
    	status TEXT,
    	response_time TEXT,
    	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	return err
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

	// checking the url
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	var status string

	if err != nil {
		if _, exist := downSince[url]; !exist {
			downSince[url] = time.Now()
		}
		status = "DOWN"
		sendEmail(url, time.Since(downSince[url]))

	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			if _, exist := downSince[url]; exist {
				delete(downSince, url)

			}
			status = "UP"
			delete(downSince, url)
		} else {
			status = fmt.Sprintf("Status %d", resp.StatusCode)
		}
	}

	_, err = db.Exec("INSERT INTO logs (website, status, response_time) VALUES (?, ?, ?)", url, status, elapsed.Milliseconds())
}

func startServer() {
	http.HandleFunc("/logs", getLogs)
	fmt.Println("REST API running on port : 8080")
	http.ListenAndServe(":8080", nil)
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	query := `
	        SELECT website, status, response_time, timestamp
			FROM logs
			ORDER BY timestamp DESC
			LIMIT 10
	`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type Log struct {
		Website      string `json:"website"`
		Status       string `json:"status"`
		ResponseTime string `json:"response_time"`
		Timestamp    string `json:"timestamp"`
	}

	var logs []Log

	for rows.Next() {
		var log Log
		rows.Scan(&log.Website, &log.Status, &log.ResponseTime, &log.Timestamp)
		logs = append(logs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func main() {
	err := loadConfig()

	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	err = initDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	defer db.Close()

	go startServer()

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
