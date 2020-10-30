package config

import (
	"fmt"
	"proxy-forward/pkg/logging"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// RuntimeViper runtime config
var RuntimeViper *viper.Viper

func init() {
	RuntimeViper = viper.New()
	RuntimeViper.SetConfigType("yaml")
	RuntimeViper.SetConfigName("cfg")         // name of config file (without extension)
	RuntimeViper.AddConfigPath("/etc/proxy/") // path to look for the config file in
	RuntimeViper.AddConfigPath("./config/")   // optionally lok for config in the working directory
	if err := RuntimeViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	RuntimeViper.WatchConfig()
	RuntimeViper.OnConfigChange(func(e fsnotify.Event) {
		logging.Log.Debugf("config file changed: %s", e.Name)
	})
}
