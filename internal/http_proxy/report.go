package http_proxy

import (
	"fmt"
	"net/http"
	"proxy-forward/config"
	"proxy-forward/internal/service/proxy_ip_service"
	"proxy-forward/pkg/e"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/gin-gonic/gin"
)

var (
	SUCCESS = fmt.Sprintf(`{"code": 0, "msg": "success", "data": {}}`)
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/", indexHandler)
	v1 := r.Group("/v1", gin.BasicAuth(gin.Accounts{
		config.RuntimeViper.GetString("server.http_authorized_username"): config.RuntimeViper.GetString("server.http_authorized_password"),
	}))
	{
		v1.POST("/destroy_ip", reportHandler)
	}
	return r
}

type reportForm struct {
	PiID int `form:"pi_id"`
}

/* 用作上报 ip 不可用 作为 内部 缓存sock 链接释放的接口 */
func reportHandler(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form reportForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	proxyIPService := proxy_ip_service.ProxyIP{ID: form.PiID}
	proxyIPService.DelteCache()

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

func indexHandler(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	appG.C.Header("Content-Type", "text/html; charset=utf-8")
	appG.C.String(200, `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title></title>
</head>
<body>
</body>
</html>
	`)
	return
}
