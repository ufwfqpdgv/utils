package utils

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ufwfqpdgv/log"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
)

func InitDB(cfg DB, outputlog bool) (db *xorm.Engine) {
	log.Debug(NowFunc())
	defer log.Debug(NowFunc() + " end")
	log.Debugf("outputlog:%v", outputlog)

	db = &xorm.Engine{}
	var connectStr string
	if cfg.Type == "mssql" {
		connectStr = fmt.Sprintf("user id=%s;password=%s;server=%s;port%d;database=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db_name)
	} else {
		connectStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db_name)
	}

	var err error
	db, err = xorm.NewEngine(cfg.Type, connectStr)
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	db.SetMapper(core.GonicMapper{})
	if outputlog {
		db.ShowSQL(true)
	} else {
		db.ShowSQL(false)
	}
	db.SetMaxIdleConns(cfg.Max_idle_conns)
	db.SetMaxOpenConns(cfg.Max_conns)
	db.SetConnMaxLifetime(time.Duration(cfg.Conn_max_lifetime) * time.Second) //这个参数好像更不释放 socket，在大并发的时候会疯狂暴涨--网上文章，使用DB.SetConnMaxLifetime(time.Second)设置连接最大复用时间，3~10秒即可。orm基本上都有相关的设置

	if outputlog {
		exist, err := pathExists(cfg.Log_path)
		if err != nil {
			log.Panic(err)
		}
		if !exist {
			err = os.Mkdir(cfg.Log_path, os.ModePerm)
			if err != nil {
				log.Panic(err)
			}
		}
		pathFileName := cfg.Log_path + cfg.Log_name
		file, err := os.Open(pathFileName)
		if err != nil && os.IsNotExist(err) {
			file, err = os.Create(pathFileName)
			if err != nil {
				log.Panic(err)
			}
		} else {
			file, err = os.OpenFile(pathFileName, os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Panic(err)
			}
		}
		db.SetLogger(xorm.NewSimpleLogger(file))
	}

	return
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
