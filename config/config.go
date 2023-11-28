package config

import (
	"fmt"
	Redis2 "github.com/go-redis/redis"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Con = &Config{}

type Config struct {
	DB *DbConfig `yaml:"db"`
}
type DbConfig struct {
	Mysql  *MysqlConfig `yaml:"mysql"`
	Redis  *Redis       `yaml:"redis"`
	Logger *Logger      `yaml:"logger"`
}
type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
	Port     string `yaml:"port"`
	Address  string `yaml:"address"`
}
type Redis struct {
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Address  string `yaml:"address"`
	DbName   int    `yaml:"dbName"`
}
type Logger struct {
	IsDebug bool `yaml:"isDebug"`
}

func (d *DbConfig) Dns(flag bool) string {

	if flag {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/", d.Mysql.Username, d.Mysql.Password, d.Mysql.Address, d.Mysql.Port)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.Mysql.Username, d.Mysql.Password, d.Mysql.Address, d.Mysql.Port, d.Mysql.DbName)
}

func InitConfig() {
	cfg := pflag.String("path", "./config/dev.yaml", "配置文件路径")
	isDebug := pflag.Bool("isdebug", true, "开始debug日志")
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfg)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
		return
	}
	err = viper.Unmarshal(Con)
	if err != nil {
		panic(err)
		return
	}
	Con.DB.Logger.IsDebug = *isDebug

}
func NewDb() (*gorm.DB, *Redis2.Client, error) {
	fmt.Println(Con.DB.Dns(false))
	db, err := gorm.Open(mysql.Open(Con.DB.Dns(false)))
	if err != nil {
		return nil, nil, err
	}
	Client := Redis2.NewClient(&Redis2.Options{
		Addr:     fmt.Sprintf("%s:%s", Con.DB.Redis.Address, Con.DB.Redis.Port),
		Password: Con.DB.Redis.Password,
		DB:       Con.DB.Redis.DbName,
	})
	err = Client.Ping().Err()
	if err != nil {
		return nil, nil, err
	}
	return db, Client, err
}
