# Website Uptime Monitoring Tool

A simple Go-based Uptime Monitoring Tool that checks the availability of websites, logs their statuses, and sends email alerts when a site is down. The tool uses SQLite to store logs and exposes a REST API to retrieve the latest status information.

## Features

- Monitors multiple websites' uptime status.
- Logs website status and response times in an SQLite database.
- Sends email alerts when a website is down.
- Exposes a simple REST API to view the latest 10 logs.
- Configurable check interval and email settings.

## Requirements

- Go (version 1.18+)
- SQLite3

## Installation

### Option 1: Install Locally

A. Clone this repository:
```bash
git clone https://github.com/MurtadaAhmed/Website-Monitoring-Tool-with-Go.git
cd Website-Monitoring-Tool-with-Go
```

B. Install dependencies:
```shell
go get github.com/mattn/go-sqlite3
go get gopkg.in/yaml.v3
```

C. Update `config.yaml` file with the websites, check interval and the email server:
```yaml
websites:
  - "https://example1.com"
  - "https://example2.com"
check_interval: 1m
email:
  smtpServer: "smtp.example.com"
  smtpPort: "587"
  sender: "your-email@example.com"
  password: "your-email-password"
  receiver: "recipient-email@example.com"
```
`websites:` List the websites you want to monitor.

`check_interval:` Set the frequency of website checks (e.g., "30s", "1m").

`email:` Provide SMTP details for sending alerts when a website is down.

D. Run the app:
```shell
go run main.go
```
The application will:
- Monitor the websites listed in the configuration file.
- Log the status and response time of each website.
- Send an email alert if any website is down.
- Expose a REST API at http://localhost:8080/logs to retrieve the latest 10 logs. The response will return the most recent 10 logs, in JSON format:
```json
[
  {
    "website": "https://example1.com",
    "status": "UP",
    "response_time": "120",
    "timestamp": "2025-03-09T12:34:56"
  },
  {
    "website": "https://example2.com",
    "status": "DOWN",
    "response_time": "0",
    "timestamp": "2025-03-09T12:35:01"
  }
]
```

## Option 2: Run with Docker

A. Update `config.yaml` file with the websites, check interval and the email server:
```yaml
websites:
  - "https://example1.com"
  - "https://example2.com"
check_interval: 1m
email:
  smtpServer: "smtp.example.com"
  smtpPort: "587"
  sender: "your-email@example.com"
  password: "your-email-password"
  receiver: "recipient-email@example.com"
```
`websites:` List the websites you want to monitor.

`check_interval:` Set the frequency of website checks (e.g., "30s", "1m").

`email:` Provide SMTP details for sending alerts when a website is down.

B. Build the Docker image:
```shell
docker build -t website-monitoring-tool-with-go .
```
C. Run the Docker container:
```shell
docker run -d -p 8080:8080 --name website-monitoring-tool website-monitoring-tool-with-go
```
This will start the monitoring tool in a Docker container, with the REST API exposed on port 8080.


## Email Alerts:
If any monitored website is down, the application will send an email to the specified receiver. The email will contain the website's URL and the duration of downtime.

Example Email:
```yaml
Subject: Website Down Alert

The website https://example2.com is down. Sun Mar  9 12:35:01 UTC 2025
Down since: 3m45s
```

## Logging
Website uptime logs are stored in an SQLite database (monitor.db). The logs include:
- Website URL 
- Status (UP/DOWN or HTTP status code)
- Response time in milliseconds 
- Timestamp of the check

## License
This project is licensed under the MIT License.