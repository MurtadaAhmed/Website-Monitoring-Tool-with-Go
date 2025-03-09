# Website-Monitoring-Tool-with-Go

**Under Development**

```shell 
CGO_ENABLED=1 go run main.go
```

API

http://localhost:8080/logs


Docker 

```shell
docker build -t website-monitoring-tool-with-go .
docker run -d -p 8080:8080 --name website-monitoring-tool website-monitoring-tool-with-go
```