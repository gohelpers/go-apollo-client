package apolloclient

import (
	"errors"
	"fmt"
	"github.com/philchia/agollo/v4"
	"github.com/zeromicro/go-zero/core/conf"
	"strings"
	"time"
)

// apollo客户端配置文件
type ApolloClientConf struct {
	Server    string `json:"server"`
	Cluster   string `json:"cluster"`
	Appid     string `json:"appid"`
	Namespace string `json:"namespace"`
	Secret    string `json:"secret"`
	CacheDir  string `json:"cacheDir"`
}

/**
 * @Description: 启动apollo
 * @Author: LiuQHui
 * @Param conf
 * @Param appConf
 * @Return error
 * @Date  2022-10-11 17:56:59
**/
func StartApollo(c ApolloClientConf, appConf interface{}) error {
	namespaces := strings.Split(c.Namespace, ",")
	cf := &agollo.Conf{
		AppID:           c.Appid,
		Cluster:         c.Cluster,
		NameSpaceNames:  namespaces,
		CacheDir:        c.CacheDir,
		MetaAddr:        c.Server,
		AccesskeySecret: c.Secret,
	}
	err := agollo.Start(cf, agollo.SkipLocalCache())
	if err != nil {
		return errors.New(fmt.Sprintf("start apollo error %v", err))
	}
	// 解析namespace内容
	content := agollo.GetContent(agollo.WithNamespace(namespaces[1]))
	_ = conf.LoadFromYamlBytes([]byte(content), appConf)
	listenApollo(c, appConf)
	return nil
}

/**
 * @Description: 监听apollo变更
 * @Author: LiuQHui
 * @Param appConf
 * @Date  2022-10-11 17:44:38
**/
func listenApollo(c ApolloClientConf, appConf interface{}) {
	go func() {
		namespaces := strings.Split(c.Namespace, ",")
		agollo.OnUpdate(func(changeEvent *agollo.ChangeEvent) {
			// 监听配置变更
			fmt.Println("event update:", agollo.GetContent(agollo.WithNamespace(namespaces[1])))
			content := agollo.GetContent(agollo.WithNamespace(namespaces[1]))
			_ = conf.LoadFromYamlBytes([]byte(content), appConf)
			fmt.Printf("appConf:%#v", appConf)
		})
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
			default:

			}
		}
	}()
}
