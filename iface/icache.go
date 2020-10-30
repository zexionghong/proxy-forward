package iface

import "net/http"

type CachePool interface {
	Get(uri string) Cache
	Delete(uri string)
	CheckAndStore(uri string, req *http.Request, resp *http.Response)
}

type Cache interface {
	Verify() bool
	WriteTo(rw http.ResponseWriter) (int, error)
}
