package nacos4viper

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	crypt "github.com/sagikazarmark/crypt/config"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net/url"
	"reflect"
)

var nacosConfigParam vo.ConfigParam

type nacosConfigManager struct {
	config_client.IConfigClient
}

func (c nacosConfigManager) Get(key string) ([]byte, error) {
	content, err := c.GetConfig(nacosConfigParam)
	return []byte(content), err
}
func (c nacosConfigManager) List(key string) (crypt.KVPairs, error) {
	return crypt.KVPairs{}, nil
}
func (c nacosConfigManager) Set(key string, value []byte) error {
	return nil
}
func (c nacosConfigManager) Watch(key string, stop chan bool) <-chan *crypt.Response {
	resp := make(chan *crypt.Response)

	if nacosConfigParam.OnChange == nil {
		nacosConfigParam.OnChange = func(namespace, group, dataId, data string) {
			resp <- &crypt.Response{
				Value: []byte(data),
				Error: nil,
			}
		}
	}
	err := c.ListenConfig(nacosConfigParam)
	if err != nil {
		return nil
	}

	go func() {
		for {
			select {
			case <-stop:
				_ = c.CancelListenConfig(nacosConfigParam)
				return
			}
		}
	}()

	return resp
}

func getConfigManager(rp viper.RemoteProvider) (crypt.ConfigManager, error) {
	urlParse, err := url.Parse(rp.Endpoint())
	if err != nil {
		return nil, errors.New("invalid endpoint")
	}

	var port uint64
	if urlParse.Port() == "0" || urlParse.Port() == "" {
		if urlParse.Scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	} else {
		port = cast.ToUint64(urlParse.Port())
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			urlParse.Host,
			cast.ToUint64(port),
			constant.WithScheme(urlParse.Scheme),
			constant.WithContextPath(rp.Path())),
	}
	cc := constant.ClientConfig{}

	if rp.SecretKeyring() != "" {
		if params, err := url.ParseQuery(rp.SecretKeyring()); err == nil && params != nil {
			ccRef := reflect.ValueOf(&cc).Elem()
			ncpRef := reflect.ValueOf(&nacosConfigParam).Elem()
			for name, val := range params {
				if f := ccRef.FieldByName(name); f.IsValid() {
					if f.CanSet() {
						switch f.Kind() {
						case reflect.String:
							f.SetString(val[0])
						case reflect.Bool:
							f.SetBool(cast.ToBool(val[0]))
						case reflect.Int:
							f.SetInt(cast.ToInt64(val[0]))
						}
					}
				} else if f := ncpRef.FieldByName(name); f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
					f.SetString(val[0])
				}
			}
			if cc.NamespaceId == "" {
				cc.NamespaceId = "public"
			}
		}
	} else {
		return nil, errors.New("SecretKeyring parameter cannot be empty")
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	return nacosConfigManager{client}, err
}

func init() {
	nacosConfigParam = vo.ConfigParam{DataId: "config", Group: "DEFAULT_GROUP"}
}
