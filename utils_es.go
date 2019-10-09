package utils

import (
	"net"
	"net/http"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

func InitEsClient(es ES) (client *elasticsearch.Client, err error) {
	Debug(NowFunc())

	for _, v := range es.Addr_arr_arr {
		cfg := elasticsearch.Config{
			Addresses: v,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConnsPerHost:   400,
				MaxIdleConns:          400,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 10 * time.Second,
			},
		}
		client, err = elasticsearch.NewClient(cfg)
		if err != nil {
			Error(err)
			continue
		}

		var res *esapi.Response
		res, err = client.Ping()
		if err != nil {
			Error(err)
			continue
		}
		defer res.Body.Close()
		if res.IsError() {
			var e map[string]interface{}
			err = Json.NewDecoder(res.Body).Decode(&e)
			if err != nil {
				Error(err)
				continue
			}
			Errorf("[%v] %v: %v",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			continue
		}

		// 已找到正常可用的ES服务器则可退出了
		if err == nil {
			break
		}
	}

	return
}
