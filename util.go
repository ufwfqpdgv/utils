package utils

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ufwfqpdgv/log"
	"github.com/ufwfqpdgv/samh_common_lib"
	"github.com/viant/toolbox"
)

var (
	Json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func Struct2Map(in interface{}) (out map[string]string) {
	m := make(map[string]interface{})
	err := toolbox.NewConverter("", "json").AssignConverted(&m, in)
	if err != nil {
		log.Error(err)
	}
	out = make(map[string]string)
	for k, v := range m {
		out[k] = fmt.Sprintf("%v", v)
	}
	return
}

func Decimal(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func DealSystemDiffData(c *gin.Context, to interface{}, from interface{}) {
	if strings.Contains(c.Request.UserAgent(), "iPhone") ||
		strings.Contains(c.Request.UserAgent(), "iPad") ||
		strings.Contains(c.Request.UserAgent(), "iOS") {
		//do nothing
	} else {
		copier.Copy(to, from)
	}
}

func IntArrayToString(a []int64, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
}

func IsIntArrayContain(arr []int64, v int64) bool {
	for _, v2 := range arr {
		if v == v2 {
			return true
		}
	}
	return false
}

func RemoveDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func IsIos(c *gin.Context) bool {
	if strings.Contains(c.Request.UserAgent(), "iPhone") ||
		strings.Contains(c.Request.UserAgent(), "iPad") ||
		strings.Contains(c.Request.UserAgent(), "iOS") {
		return true
	} else {
		return false
	}
}

func IsAndroidOldVersion(c *gin.Context) bool {
	// FIXME 抓包看到的安卓发出的收到的header 里的 key、key 数量都对不上，不懂是 nginx 转发问题还是此框架的解析问题
	// 因不替换老版本接口，故会全是新的版本，直接当新的用先
	return false

	if strings.Contains(c.GetHeader("os"), "android") {
		versionArr := strings.Split(c.GetHeader("app_version"), ".")
		versionS := fmt.Sprintf("%03v%03v%03v", versionArr[0], versionArr[1], versionArr[2])
		version, _ := strconv.Atoi(versionS)
		if version >= 2001000 { // 旧配置 version 2.1.0->2001000
			return false
		}
	}

	return true
}

func NewConfigWatcher(env string, init func()) {
	Info(NowFunc())
	defer Info(NowFunc() + " end")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					Error(NowFuncError())
					return
				}

				cmdStr := fmt.Sprintf("cd config/%v && ls *.*", env)
				cmd := exec.Command("sh", "-c", cmdStr)
				out, err := cmd.Output()
				if err != nil {
					var stderr bytes.Buffer
					cmd.Stderr = &stderr
					Errorf("%v\n%v", err, stderr)
					// golang 执行 shell 语句时会莫名其秒的报错，实际上结果是对的，如这里的改 toml 文件的时候，在配置项最后加上多余的空格这里也是err并提示"ls: 无法访问dev_config.toml~: 没有那个文件或目录"
					continue
				}
				notifyFileArr := strings.Split(string(out), "\n")

				for _, v := range notifyFileArr {
					if event.Name == fmt.Sprintf("config/%v/%v", env, v) &&
						event.Op&fsnotify.Write == fsnotify.Write {
						Info("modified file:", event.Name)
						init()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					Error(NowFuncError())
					return
				}
				Error(err)
			}
		}
	}()

	err = watcher.Add(fmt.Sprintf("config/%v", env))
	if err != nil {
		Panic(err)
	}
	<-done
}

func ArrToSet(arr interface{}) (set mapset.Set, retCode samh_common_lib.SamhResponseCode) {
	Debug(NowFunc())
	defer Debug(NowFunc() + " end")

	retCode = samh_common_lib.SamhResponseCode_Succ

	set = mapset.NewSet()
	switch arr.(type) {
	case []int:
		for _, v := range arr.([]int) {
			set.Add(v)
		}
	case []int64:
		for _, v := range arr.([]int64) {
			set.Add(v)
		}
	case []string:
		for _, v := range arr.([]string) {
			set.Add(v)
		}
	default:
		retCode = samh_common_lib.SamhResponseCode_Param_Invalid
		return
	}

	return
}

func GetRequestUrl(url string) (urlRsp string) {
	index := strings.Index(url, "?")
	if index >= 0 {
		urlRsp = url[:index]
	} else {
		urlRsp = url
	}

	return
}

func CheckUrlTimeout(rqUrl string, rqTimeout int) (rspUrl string, rspTimeout int, rspErr error) {
	if rqUrl == "" || rqTimeout == 0 {
		rspErr = errors.New("url或 timeOut为空")
		return
	}
	rspUrl, rspTimeout = rqUrl, rqTimeout

	return
}

func CheckServerConnect(rqUrl string) (rspErr error) {
	errC := make(chan error, 1)

	go func(rqUrlTemp string, errCTemp chan error) {
		s1 := strings.Split(rqUrl, "//")
		s2 := strings.Split(s1[1], ":")
		host, port := s2[0], s2[1]
		cmdStr := fmt.Sprintf("telnet %v %v", host, port)
		cmd := exec.Command("sh", "-c", cmdStr)
		out, err := cmd.Output()
		if err != nil {
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			Errorf("%v\n%v", err, stderr)
			// golang 执行 shell 语句时会莫名其秒的报错，实际上结果是对的
		}
		exist := strings.Contains(string(out), "Connected")
		if !exist {
			errCTemp <- errors.New(fmt.Sprintf("can't connect %v", rqUrlTemp))
			return
		}
		errCTemp <- nil
	}(rqUrl, errC)

	select {
	case err := <-errC:
		rspErr = err
	case <-time.After(time.Second * 5):
		rspErr = errors.New(fmt.Sprintf("connect time out,url:%v", rqUrl))
	}

	return
}

func VersionStr2Int(version string) (versionRsp int, err error) {
	versionArr := strings.Split(version, ".")
	if len(versionArr) != 3 {
		err = errors.New("参数错误")
		return
	}
	versionS := fmt.Sprintf("%03v%03v%03v", versionArr[0], versionArr[1], versionArr[2]) // 2.1.0->2001000
	versionRsp, err = strconv.Atoi(versionS)
	if err != nil {
		return
	}

	return
}

func NowFunc() string {
	pc, _, _, _ := runtime.Caller(1)
	return "NowFunc:" + runtime.FuncForPC(pc).Name() + " "
}

func NowFuncError() string {
	pc, _, _, _ := runtime.Caller(1)
	return "NowFunc:" + runtime.FuncForPC(pc).Name() + " Error "
}
