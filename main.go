// package main
//
// import (
// 	"fmt"
// 	"log"
// 	"proxy-forward/pkg/utils"
// 	"sync"
// 	"time"
//
// 	"github.com/parnurzeal/gorequest"
// )
//
// var proxyIP = "http://47.113.114.143"
// var (
// 	success = 0
// 	fail    = 0
// )
//
// func main() {
// 	ch := make(chan int, 50)
// 	var wg sync.WaitGroup
// 	startTime := time.Now()
//
// 	for i := 10000; i < 15000; i++ {
// 		ch <- 0
// 		wg.Add(1)
// 		go Get(i, ch, &wg)
// 	}
// 	wg.Wait()
// 	fmt.Printf("sueecss: %d, fail: %d", success, fail)
// 	endTime := time.Now()
// 	fmt.Printf("cost time: %.2f", endTime.Sub(startTime).Seconds())
// 	log.Println(utils.InetAtoN("47.113.114.143"))
// }
//
// func Get(port int, ch chan int, wg *sync.WaitGroup) {
// 	defer func() {
// 		wg.Done()
// 		<-ch
// 	}()
// 	request := gorequest.New().Timeout(time.Duration(5 * time.Second)).Proxy(fmt.Sprintf("%s:%d", proxyIP, port))
// 	// _, body, err := request.Get("http://myip.ipip.net/").End()
// 	// if err != nil {
// 	// 	fmt.Printf("%d - err: %v", port, err)
// 	// 	fail += 1
// 	// 	return
// 	// }
// 	success += 1
// 	// fmt.Printf("%d - %s", port, body)
// }
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ouqiang/goproxy"
)

type EventHandler struct{}

func (e *EventHandler) Connect(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 保存的数据可以在后面的回调方法中获取
	ctx.Data["req_id"] = "uuid"

	// 禁止访问某个域名
	if strings.Contains(ctx.Req.URL.Host, "example.com") {
		rw.WriteHeader(http.StatusForbidden)
		ctx.Abort()
		return
	}
}

func (e *EventHandler) Auth(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 身份验证
}

func (e *EventHandler) BeforeRequest(ctx *goproxy.Context) {
	// 修改header
	ctx.Req.Header.Add("X-Request-Id", ctx.Data["req_id"].(string))
	// 设置X-Forwarded-For
	if clientIP, _, err := net.SplitHostPort(ctx.Req.RemoteAddr); err == nil {
		if prior, ok := ctx.Req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		ctx.Req.Header.Set("X-Forwarded-For", clientIP)
	}
	// 读取Body
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		// 错误处理
		return
	}
	// Request.Body只能读取一次, 读取后必须再放回去
	// Response.Body同理
	ctx.Req.Body = ioutil.NopCloser(bytes.NewReader(body))

}

func (e *EventHandler) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	// 修改response
}

// 设置上级代理
func (e *EventHandler) ParentProxy(req *http.Request) (*url.URL, error) {
	a, err := url.Parse("http://localhost:8123")
	return a, err
}

func (e *EventHandler) Finish(ctx *goproxy.Context) {
	fmt.Printf("请求结束 URL:%s\n", ctx.Req.URL)
}

// 记录错误日志
func (e *EventHandler) ErrorLog(err error) {
	log.Println(err)
}

func main() {
	proxy := goproxy.New(goproxy.WithDelegate(&EventHandler{}))
	server := &http.Server{
		Addr:         ":8080",
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
