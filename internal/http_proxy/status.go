package http_proxy

var (
	HTTP200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")
	HTTP407 = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Access to internal site\"\r\n\r\n") //（需要代理授权） 此状态代码与 401（未授权）类似，但指定请求者应当授权使用代理。
	HTTP502 = []byte("HTTP/1.1 502 Bad Gateway\r\n\r\n")                                                                                 // （错误网关） 服务器作为网关或代理，从上游服务器收到无效响应。
	HTTP503 = []byte("HTTP/1.1 503 Service Unavailable\r\n\r\n")
	HTTP504 = []byte("HTTP/1.1 504 Gateway Timeout\r\n\r\n") // （网关超时）  服务器作为网关或代理，但是没有及时从上游服务器收到请求。
)
