package utils

import (
	"github.com/jinzhu/configor"
)

var _cfg *Config

func ConfigInit(configFilePath string) {
	_cfg = &Config{}
	err := configor.Load(_cfg, configFilePath)
	if err != nil {
		panic(err)
	}
}

func Init2(configFilePath string, cfg interface{}) {
	err := configor.Load(cfg, configFilePath)
	if err != nil {
		panic(err)
	}
}

func ConfigInstance() (cfg *Config) {
	return _cfg
}

type Config struct {
	Base_info struct {
		Version string
		Name    string
		Port    int
		App_id  int
	}

	Log_info_item   Log_info
	Internal_server map[string]Internal_serverStruct

	DB_whole_item      DBWhole
	Redis_item         Redis
	Redis_cluster_item RedisCluster

	Web struct {
		Http_request_timeout int
	}

	Sentry_dsn_item Sentry_dsn

	ES_item ES
}

type Log_info struct {
	Level            string
	Encoding         string
	Stdout           bool
	Development_mode bool
	Path_filename    string
	Max_size         int
	Max_backups      int
	Max_age          int
	Compress         bool
}

type Internal_serverStruct struct {
	Url      string
	Time_out int
}

type DBWhole struct {
	Is_use     bool
	Output_log bool
	DB_arr     map[string]DB
}

type DB struct {
	Type              string
	Host              string
	Port              int
	User              string
	Password          string
	Db_name           string
	Max_conns         int
	Max_idle_conns    int
	Conn_max_lifetime int
	Time_out          int
	Log_path          string
	Log_name          string
	//
	Table_name map[string]string
}

type Redis struct {
	Is_use      bool
	Network     string
	Addr        string
	Password    string
	Max_retries int
	Pool_size   int
	Prefix      string
	Time_out    int
}

type RedisCluster struct {
	Is_use          bool
	Master_addr_arr []string
	Slave_addr_arr  []string
	Password        string
	Max_retries     int
	Pool_size       int
	Prefix          string
	Time_out        int
}

type Sentry_dsn struct {
	Is_use bool
	Url    string
}

type ES struct {
	Is_use       bool
	Addr_arr_arr [][]string
}
