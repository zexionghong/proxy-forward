package handler

import (
	"net/http"
	"proxy-forward/internal/models"
)

func init() {
}

func (hs *HandlerServer) LoadTraveling(userToken *models.UserToken) (*http.Transport, error) {
	return nil, nil
}

func (hs *HandlerServer) Done(rw http.ResponseWriter, req *http.Request) {
	// TODO: 流量统计
}
