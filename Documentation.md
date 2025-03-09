****1. sendEmail()****
- This function is used to send email notification if a website is down.
- It is called in the checkWebsite() function.
- (message) combines the subject and the body of the email.
- (auth) initialize the authentication using the Google app password.
- (err) this will send the email and store the error if any in the variable

****2. checkWebsite()****
- This function is used to check the status of a website.
- It creates a log file log.txt (if it doesn't exist) and write the logs to it.
- This function is called in the main() function, and it is run using goroutines (go checkWebsite(url, &wg)) and also has
the wait group wg *sync.WaitGroup to wait for the goroutines to finish.
- This line:
```shell
file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```
os.O_APPEND: used to append to the file, os.O_CREATE: to create the file if it doesn't exist,
os.O_WRONLY: to open the file in readonly mode
- (logMessage) variable store all the log messages for the website statuses
- (file) variable refers to (log.txt) file and stores all the log messages (logMessage) 
****3. main()****
- This function loops through the websites (var websites), increase the WaitGroup counter wg.Add(1), and call
the checkWebsite() function through goroutines (go checkWebsite(url, &wg))