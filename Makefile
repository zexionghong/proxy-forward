build proxy-forward:
	GOOS=linux GOARCH=amd64 go build -o proxy-forward cmd/main.go
