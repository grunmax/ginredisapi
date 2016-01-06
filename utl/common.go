package utl

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

var errorlog *os.File
var Logger *log.Logger

type Config struct {
	RedisUrl       string
	RedisPassword  string
	HttpUrl        string
	MaxConnections int
}

type MyErr struct {
	Msg string
}

//if _, ok := err.(MyErr); ok {
//  // Handle MyError
//} else {
//  // Handle all other error types
//}

func (e MyErr) Error() string {
	return e.Msg
}

func Err(userMessage string, e interface{}) {
	if e != nil {
		s := fmt.Sprintf("ERROR:%s  %v\n", userMessage, e)
		fmt.Printf(s)
		Logger.Panicf(s)
	}
}

func Log(userMessage string, v interface{}) {
	if v != nil {
		s := fmt.Sprintf(":%s  %v\n", userMessage, v)
		fmt.Printf(s)
		Logger.Printf(s)
	}
}

func InitLog() {
	errorlog, e := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if e != nil {
		fmt.Println(e)
	}
	Logger = log.New(errorlog, "app: ", log.Lshortfile|log.LstdFlags)
}

func ReadConfig() *Config {
	iniFile := "app.ini"
	config := new(Config)
	cfg, err := ini.Load([]byte(""), iniFile)
	Err("no config file", err)
	config.RedisUrl = cfg.Section("redis").Key("url").String()
	config.RedisPassword = cfg.Section("redis").Key("password").String()
	config.MaxConnections, err = cfg.Section("redis").Key("maxconnections").Int()
	Err("Wrong value for maxconnections", err)
	config.HttpUrl = cfg.Section("gin").Key("url").String()
	return config
}
