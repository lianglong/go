package nacos4viper

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestNacosConfig(t *testing.T) {
	remoteViper := viper.New()
	SetViper(remoteViper)
	//or  remoteViper := nacos4viper.New() and nacos4viper.SetViper(remoteViper)
	remoteViper.SetConfigName("test")
	remoteViper.SetConfigType("yaml")
	remoteViper.AddConfigPath("./")
	if err := remoteViper.ReadInConfig(); err != nil {
		t.Error(err)
	}
	if configCenter := remoteViper.Sub("setting.config_center"); configCenter != nil && configCenter.GetBool("enabled") {
		if err := remoteViper.AddSecureRemoteProvider(
			configCenter.GetString("provider"),
			configCenter.GetString("endpoint"),
			configCenter.GetString("path"),
			configCenter.GetString("client_param"),
		); err != nil {
			t.Error(err)
		}

		if err := remoteViper.ReadRemoteConfig(); err != nil {
			t.Error(err)
		}

		fmt.Println("setting.database.default.host[remote] ", remoteViper.GetString("setting.database.default.host"))

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill)

		go remoteViper.WatchRemoteConfigOnChannel()
		go func() {
			for {
				time.Sleep(time.Second * 5)
				fmt.Println("setting.database.default.host[watch] ", remoteViper.GetString("setting.database.default.host"))
			}
		}()
		s := <-c
		fmt.Println("stop,signal:", s)
	}

}
