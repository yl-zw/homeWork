package config

import (
	"fmt"
	Redis2 "github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Con = &Config{}

type Config struct {
	DB *DbConfig `yaml:"DB"`
}
type DbConfig struct {
	Mysql  *MysqlConfig `yaml:"mysql"`
	Redis  *Redis       `yaml:"redis"`
	ALiYun *AliYun      `yaml:"ALiYun"`
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
type AliYun struct {
	KeyID  string `yaml:"key_id"`
	Secret string `yaml:"secret"`
}

func (d *DbConfig) Dns(flag bool) string {

	if flag {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/", d.Mysql.Username, d.Mysql.Password, d.Mysql.Address, d.Mysql.Port)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.Mysql.Username, d.Mysql.Password, d.Mysql.Address, d.Mysql.Port, d.Mysql.DbName)
}

func InitConfig(configPath string) {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(Con)
	if err != nil {
		fmt.Println(err)
		return
	}
	if Con.DB.Mysql == nil || Con.DB.Redis == nil {
		fmt.Println("读取配置文件失败")
		return
	}

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
