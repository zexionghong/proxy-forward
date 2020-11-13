package app

import (
	"proxy-forward/pkg/logging"

	"github.com/astaxie/beego/validation"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Log.Infof("%s,%s", err.Key, err.Message)
	}
	return
}
