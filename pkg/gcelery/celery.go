package gcelery

import (
	"proxy-forward/config"

	"github.com/gocelery/gocelery"
)

var CeleryBroker *gocelery.RedisCeleryBroker
var CeleryBackend *gocelery.RedisCeleryBackend
var CeleryCli *gocelery.CeleryClient

var (
	ForwardTaskName = "tasks.forward_data_task"
)

// Setup initialize the gocelery instance
func Setup() error {
	var err error
	CeleryBroker = gocelery.NewRedisCeleryBroker(config.RuntimeViper.GetString("celery.broker"))
	CeleryBackend := gocelery.NewRedisCeleryBackend(config.RuntimeViper.GetString("celery.backend"))
	CeleryCli, err = gocelery.NewCeleryClient(CeleryBroker, CeleryBackend, 1)
	if err != nil {
		return err
	}
	return nil
}

// Send Task {"remote_addr": "", "usage": "", "user_token_id": 1, "type": "req/resp"}
func SendForwardDataTask(data interface{}) error {
	if config.RuntimeViper.GetBool("celery.status") == true {
		_, err := CeleryCli.Delay(ForwardTaskName, data)
		return err
	}
	return nil
}
