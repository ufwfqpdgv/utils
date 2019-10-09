package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	resty "gopkg.in/resty.v1"
)

func init() {
	clientMap["comic"], restyClientMap["comic"] = newClinet()
	clientMap["vip"], restyClientMap["vip"] = newClinet()
	clientMap["operation"], restyClientMap["operation"] = newClinet()
	clientMap["user"], restyClientMap["user"] = newClinet()
	clientMap["pay"], restyClientMap["pay"] = newClinet()
	clientMap["other"], restyClientMap["other"] = newClinet()
}

var (
	clientMap      = make(map[string]*http.Client)
	restyClientMap = make(map[string]*resty.Client)
)

// http 连接池
func newClinet() (client *http.Client, restyClient *resty.Client) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   400,
		MaxIdleConns:          400,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	client = &http.Client{
		Timeout:   time.Duration(2) * time.Second,
		Transport: transport,
	}
	restyClient = resty.NewWithClient(client)

	return
}

func getClient(url string) (key string) {
	if strings.Contains(url, "/internal/api/v1/comics/") {
		key = "comic"
	} else if strings.Contains(url, "/internal/api/v1/vip/") {
		key = "vip"
	} else if strings.Contains(url, "/internal/api/v1/operation/") {
		key = "operation"
	} else if strings.Contains(url, "/userapi/internal/") {
		key = "user"
	} else if strings.Contains(url, "/app_api/pay/") {
		key = "pay"
	} else {
		key = "other"
	}

	return
}

func HttpGet(url string, rq interface{}, rsp interface{}, timeout int, headers map[string]string, retryCount int) (err error) {
	Debug(NowFunc())
	defer Debug(NowFunc() + " end")

	Infof("GET,Url:%v,Request:%+v", url, rq)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	request.Header.Add("Connection", "close")

	rqMap := Struct2Map(rq)
	requestParam := request.URL.Query()
	for k, v := range rqMap {
		requestParam.Add(k, v)
	}
	request.URL.RawQuery = requestParam.Encode()

	response, err := clientMap[getClient(url)].Do(request)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err, response)
		return
	}
	defer response.Body.Close()
	response.Header.Add("Connection", "close")

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}
	err = Json.Unmarshal(responseBytes, rsp)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}

	Infof("Response:%+v", rsp)

	return
}

func HttpPost(url string, rq interface{}, rsp interface{}, timeout int, headers map[string]string, retryCount int) (err error) {
	Debug(NowFunc())
	defer Debug(NowFunc() + " end")

	Infof("POST,Url:%v,Request:%+v", url, rq)

	bytesData, err := json.Marshal(rq)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	request.Header.Add("Connection", "close")
	request.Header.Set("Content-Type", "application/json")

	response, err := clientMap[getClient(url)].Do(request)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}
	defer response.Body.Close()
	response.Header.Add("Connection", "close")

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}
	err = Json.Unmarshal(responseBytes, rsp)
	if err != nil {
		err = errors.New("服务器错误")
		Error(err)
		return
	}

	Infof("Response:%+v", rsp)

	return
}

/* func HttpGet(url string, rq interface{}, rsp interface{}, timeout int, headers map[string]string, retryCount int) (retCode base.SamhResponseCode) {
	log.Debug(base.NowFunc())
	defer log.Debug(base.NowFunc() + " end")

	log.Infof("GET,Url:%v,Request:%+v", url, rq)
	var (
		err  error
		resp *resty.Response
	)

	retryCount = 1
	resty.SetTimeout(time.Duration(timeout) * time.Second).
		SetRetryCount(retryCount)
	if len(headers) > 0 {
		resp, err = restyClientMap[getClient(url)].R().
			SetHeaders(headers).
			SetQueryParams(Struct2Map(rq)).
			SetResult(rsp).
			Get(url)
	} else {
		resp, err = restyClientMap[getClient(url)].R().
			SetQueryParams(Struct2Map(rq)).
			SetResult(rsp).
			Get(url)
	}
	if err != nil {
		retCode = base.SamhResponseCode_ServerError
		log.Error(err, resp)
		return
	}
	retCode = base.SamhResponseCode_Succ
	log.Infof("Response:%+v", rsp)

	return
}

func HttpPost(url string, rq interface{}, rsp interface{}, timeout int, headers map[string]string, retryCount int) (retCode base.SamhResponseCode) {
	log.Debug(base.NowFunc())
	defer log.Debug(base.NowFunc() + " end")

	log.Infof("POST,Url:%v,Request:%+v", url, rq)
	var (
		err  error
		resp *resty.Response
	)

	retryCount = 1
	resty.SetTimeout(time.Duration(timeout) * time.Second).
		SetRetryCount(retryCount)
	if len(headers) > 0 {
		resp, err = restyClientMap[getClient(url)].R().
			SetHeaders(headers).
			SetBody(rq).
			SetResult(rsp).
			Post(url)
	} else {
		resp, err = restyClientMap[getClient(url)].R().
			SetBody(rq).
			SetResult(rsp).
			Post(url)
	}
	if err != nil {
		retCode = base.SamhResponseCode_ServerError
		log.Error(err, resp)
		return
	}
	retCode = base.SamhResponseCode_Succ
	log.Infof("Response:%+v", rsp)

	return
} */

//旧版本的只支持这样form的并把rq先转成json串的样式
func HttpPost2(url string, rq interface{}, rsp interface{}, timeout int) (err error) {
	Debug(NowFunc())
	defer Debug(NowFunc() + " end")

	Infof("POST2,Url:%v,Request:%+v", url, rq)
	resty.SetTimeout(time.Duration(timeout) * time.Second)
	b, err := Json.Marshal(rq)
	if err != nil {
		Error(err.Error())
		err = errors.New("参数无效")
		return
	}
	resp, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(string(b)).
		SetResult(rsp).
		Post(url)
	if err != nil {
		err = errors.New("参数无效")
		Error(err, resp)
		return
	}
	Infof("Response:%+v", rsp)

	return
}
