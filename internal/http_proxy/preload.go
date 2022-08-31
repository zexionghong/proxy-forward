package http_proxy

import (
	"errors"
	"net/http"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/bypass_service"
	"proxy-forward/pkg/logging"
	"proxy-forward/pkg/utils"
	"strings"
)

func (hs *HandlerServer) PreloadReq(userToken *models.UserToken, rw http.ResponseWriter, req *http.Request) bool {
	err := hs.preloadReq(userToken, rw, req)
	if err != nil {
		return false
	}
	return true
}

// preload req url
func (hs *HandlerServer) preloadReq(userToken *models.UserToken, rw http.ResponseWriter, req *http.Request) error {
	if userToken == nil {
		Unavailable(rw)
		return errors.New("userToken is nil")
	}
	if userToken.PsID == 0 {
		return nil
	}
	bypassService := bypass_service.Bypass{PsID: userToken.PsID}
	bypassDomain, err := bypassService.GetBypassByPsID()
	if err != nil {
		Unavailable(rw)
		return err
	}
	hostname := strings.TrimPrefix(req.URL.Hostname(), "www.")
	exists := utils.BelongsToList(hostname, bypassDomain)
	if exists {
		logging.Log.Infof("请求黑名单域名 uid:%d, hostname: %s", userToken.Uid, hostname)
		Unavailable(rw)
		return errors.New("domain is not bypass")
	}
	return nil
}
