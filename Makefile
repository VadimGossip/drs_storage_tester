build:
	GOOS=linux GOARCH=amd64 go build -o chat-server cmd/main.go