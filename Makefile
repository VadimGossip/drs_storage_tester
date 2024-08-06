build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o dbstand cmd/main.go