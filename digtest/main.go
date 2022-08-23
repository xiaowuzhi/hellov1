package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"

	"go.uber.org/dig"
	"github.com/go-ini/ini"
)

type Option struct {
	ConfigFile string `short:"c" long:"config" description:"Name of config file."`
}

type RedisConfig struct {
	IP   string
	Port int
	DB   int
}

type MySQLConfig struct {
	IP       string
	Port     int
	User     string
	Password string
	Database string
}

type Config struct {
	dig.In

	Redis *RedisConfig
	MySQL *MySQLConfig
}

func InitOption() (*Option, error) {
	var opt Option
	opt = Option{ConfigFile: "abc.text"}
	_, err := flags.Parse(&opt)
	return &opt, err
}

func InitConfig(opt *Option) (*ini.File, error) {
	cfg, err := ini.Load(opt.ConfigFile)
	return cfg, err
}

func InitRedisConfig(cfg *ini.File) (*RedisConfig, error) {
	port, err := cfg.Section("redis").Key("port").Int()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db, err := cfg.Section("redis").Key("db").Int()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &RedisConfig{
		IP:   cfg.Section("redis").Key("ip").String(),
		Port: port,
		DB:   db,
	}, nil

}

func InitMySQLConfig(cfg *ini.File) (*MySQLConfig, error) {
	port, err := cfg.Section("mysql").Key("port").Int()
	if err != nil {
		return nil, err
	}
	return &MySQLConfig{
		IP:       cfg.Section("mysql").Key("ip").String(),
		Port:     port,
		User:     cfg.Section("mysql").Key("user").String(),
		Password: cfg.Section("mysql").Key("password").String(),
		Database: cfg.Section("mysql").Key("database").String(),
	}, nil
}

func PrintInfo(config Config) {
	fmt.Println("================== redis section =====================")
	fmt.Println("redis ip:", config.Redis.IP)
	fmt.Println("redis port:", config.Redis.Port)
	fmt.Println("redis db:", config.Redis.DB)

	fmt.Println("================== mysql section =====================")
	fmt.Println("mysql ip:", config.MySQL.IP)
	fmt.Println("mysql port:", config.MySQL.Port)
	fmt.Println("mysql user:", config.MySQL.User)
	fmt.Println("mysql password:", config.MySQL.Password)
	fmt.Println("mysql db:", config.MySQL.Database)

}

type DbMgr struct {
	config Config
}

func InitMgr(c Config) *DbMgr {
	return &DbMgr{config: c}
}

func (db *DbMgr) Connect() {
	fmt.Println("~~~~~~~~~~~~~~~~~~ redis section connect ~~~~~~~~~~~~~")
	fmt.Println("redis ip:", db.config.Redis.IP)
	fmt.Println("redis port:", db.config.Redis.Port)
	fmt.Println("redis db:", db.config.Redis.DB)
}

func main() {
	container := dig.New()

	container.Provide(InitOption)
	container.Provide(InitConfig)
	container.Provide(InitMySQLConfig)
	container.Provide(InitRedisConfig)

	err := container.Invoke(PrintInfo)
	if err != nil {
		log.Fatal(err)
	}

	//1.Privide连接数据库
	container.Provide(InitMgr)
	//2获取DbMgr
	var mgr *DbMgr
	err = container.Invoke(func(m *DbMgr) {
		mgr = m
	})

	if err != nil {
		log.Fatal(err)
	}

	//3执行DbMgr的Connect方法
	err = container.Invoke(mgr.Connect)
	if err != nil {
		log.Fatal(err)
	}

}
