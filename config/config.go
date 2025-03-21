package config

import (
	"fmt"
	"os"
	"path"
	"schoolonline/launch"
	"time"

	gocfg "github.com/dsbasko/go-cfg"
)

type Config struct {
	BotKey     string `yaml:"botKey"`
	BotKeyTest string `yaml:"botKeyTest"`
	EmailKey   string `yaml:"emailKey"`
	EmailSend  string `yaml:"emailSend"`
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	Port       int    `yaml:"port"`
	DbName     string `yaml:"dbname"`
	Password   string `yaml:"password"`
	SslMode    string `yaml:"sslmode"`
	SessionKey string `yaml:"sessionKey"`
	UrlSite    string `yaml:"urlsite"`
	UrlType    string `yaml:"urltype"`
	UrlPort    int    `yaml:"urlport"`
	LinkServer string `yaml:"link_server"`
	LinkTest   string `yaml:"link_test"`
	CoderKey   string `yaml:"coderKey"`
	ConnServ   string
	Launch     string
}

var C Config

func GetConfig() {

	launch.IsLaunch()

	if err := gocfg.ReadFile(path.Join("./config.yaml"), &C); err != nil {
		fmt.Printf("failed to read config.yaml file: %v", err)
		<-time.After(10 * time.Second)
		os.Exit(3)
	}

	C.createLaunch().createConnServ().Print()
}

// генерация строки подключения
func (a *Config) createConnServ() *Config {
	a.ConnServ = fmt.Sprintf("host=%s port=%d user=%s dbname='%s' password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DbName, a.Password, a.SslMode)
	return a
}

func (a *Config) createLaunch() *Config {
	if launch.Launch == "server" {
		a.Host = "localhost"
		a.Launch = "server"
	} else {
		a.Launch = "home"
	}
	return a
}

// печать конфига
func (a *Config) Print() {
	fmt.Println("----------------------------------")
	fmt.Println("структура: Config")
	fmt.Println("Launch    :          ", a.Launch)
	fmt.Println("BotKey    :          ", a.BotKey)
	fmt.Println("EmailKey  :          ", a.EmailKey)
	fmt.Println("EmailSend :          ", a.EmailSend)
	fmt.Println("SessionKey:          ", a.SessionKey)
	fmt.Println("Host      :          ", a.Host)
	fmt.Println("User      :          ", a.User)
	fmt.Println("Port      :          ", a.Port)
	fmt.Println("DbName    :          ", a.DbName)
	fmt.Println("Password  :          ", a.Password)
	fmt.Println("SslMode   :          ", a.SslMode)
	fmt.Println("ConnServ  :          ", a.ConnServ)
	fmt.Println("UrlSite   :          ", a.UrlSite)
	fmt.Println("UrlType   :          ", a.UrlType)
	fmt.Println("----------------------------------")
	fmt.Println("")
}
