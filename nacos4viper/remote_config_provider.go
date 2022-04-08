package nacos4viper

import (
	"bytes"
	crypt "github.com/sagikazarmark/crypt/config"
	"github.com/spf13/viper"
	"io"
)

var n4viper *viper.Viper

type nacosRemoteConfigProvider struct {
}

func (rc nacosRemoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	b, err := cm.Get(rp.Path())
	r := bytes.NewReader(b)
	//Overwrite data in viper.ReadInConfig()
	if err == nil && getViper().InConfig(nacosConfigParam.DataId) {
		getViper().MergeConfig(r)
	}
	return r, err
}

func (rc nacosRemoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return rc.Get(rp)
}

func (rc nacosRemoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, nil
	}
	quit := make(chan bool)
	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *crypt.Response, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				quit <- true
				return
			case resp := <-cr:
				//Overwrite data in viper.ReadInConfig()
				if resp.Error == nil && getViper().InConfig(nacosConfigParam.DataId) {
					getViper().MergeConfig(bytes.NewReader(resp.Value))
				}

				vr <- &viper.RemoteResponse{
					Error: resp.Error,
					Value: resp.Value,
				}
			}
		}
	}(cryptoResponseCh, viperResponsCh, quitwc, quit)

	return viperResponsCh, quitwc
}

func init() {
	if ContainsString(viper.SupportedRemoteProviders, "nacos") < 0 {
		viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, "nacos")
	}
	viper.RemoteConfig = &nacosRemoteConfigProvider{}
}

func ContainsString(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

func New() *viper.Viper {
	n4viper = viper.New()
	return n4viper
}

func SetViper(v *viper.Viper) {
	n4viper = v
}

func getViper() *viper.Viper {
	if n4viper == nil {
		panic("Please use nacos4viper.New() or nacos4viper.SetViper() methods to pass context variables")
	}
	return n4viper
}
