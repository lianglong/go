Go spf13/viper的nacos远程配置实现

## 快速使用

```go
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

	if err := viper.AddSecureRemoteProvider(
		"nacos",
		"http://127.0.0.1:8848",
		"/nacos",
		"NamespaceId=nacos4viper&Username=nacos4viper&Password=nacos4viper&DataId=setting&Group=dev",
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

```


## LICENSE

**[MIT](LICENSE)**