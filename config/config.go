package config

import (
	"fmt"
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
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
	Port     string `yaml:"port"`
	Address  string `yaml:"address"`
}

func (d *DbConfig) Dns(flag bool) string {

	if flag {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/", d.Username, d.Password, d.Address, d.Port)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.Username, d.Password, d.Address, d.Port, d.DbName)
}

func init() {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(Con)
	if err != nil {
		fmt.Println(err)
		return
	}

}
func NewDb() (*gorm.DB, error) {
	fmt.Println(Con.DB.Dns(false))
	db, err := gorm.Open(mysql.Open(Con.DB.Dns(false)))
	return db, err
}
