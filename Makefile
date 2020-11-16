build proxy-forward:
	GOOS=linux GOARCH=amd64 go build -o proxy-http-forward cmd_http_proxy/main.go
