buildall: http-proxy-forward socks-proxy-forward

http-proxy-forward:
	GOOS=linux GOARCH=amd64 go build -o http-proxy-forward cmd_http_proxy/main.go

socks-proxy-forward:
	GOOS=linux GOARCH=amd64 go build -o socks-proxy-forward cmd_socks_proxy/main.go
