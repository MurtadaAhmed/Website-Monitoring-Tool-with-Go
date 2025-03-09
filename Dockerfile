FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o main
EXPOSE 8080
CMD ["./main"]