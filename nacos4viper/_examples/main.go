package main

import (
	"fmt"
	"github.com/lianglong/go/nacos4viper"
	"os"
	"os/signal"
	"time"
)

func main() {
	viper := nacos4viper.New()
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.AddSecureRemoteProvider(
		viper.GetString("setting.config_center.provider"),
		viper.GetString("setting.config_center.endpoint"),
		viper.GetString("setting.config_center.path"),
		viper.GetString("setting.config_center.client_param"),
	); err != nil {
		panic(err)
	}
	if err := viper.ReadRemoteConfig(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	go viper.WatchRemoteConfigOnChannel()

	//This coroutine is not required, it is only used as a test
	go func() {
		for {
			time.Sleep(time.Second * 5)
			fmt.Println("setting.database.default.host[watch] ", viper.GetString("setting.database.default.host"))
		}
	}()
	<-c
}
