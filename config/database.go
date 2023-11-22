package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.qunhequnhe.com/coops/toad-goclient/client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Toad struct {
	Toad ToadConfig `yaml:"toad"`
}
type ToadConfig struct {
	Appid string `yaml:"appId"`
	Token string `yaml:"token"`
	Stage string `yaml:"stage"`
}

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB(configPath *string) error {
	// 读取 YAML 文件
	yamlFile, configErr := ioutil.ReadFile(*configPath)
	if configErr != nil {
		log.Fatalf("Error reading YAML file: %v", configErr)
	}

	// 解析 YAML 数据
	var config Toad
	configErr = yaml.Unmarshal(yamlFile, &config)
	if configErr != nil {
		log.Fatalf("Error unmarshalling YAML: %v", configErr)
	}

	var mtc = client.NewMiddlewareToadClient("app", config.Toad.Appid, config.Toad.Stage, config.Toad.Token, []string{"confluence.database.host", "confluence.database.username", "confluence.database.password"})
	mtc.InitConfig()
	var err error
	host, _ := mtc.GetConfigByKey("confluence.database.host")
	username, _ := mtc.GetConfigByKey("confluence.database.username")
	password, _ := mtc.GetConfigByKey("confluence.database.password")
	fmt.Print(username + ":" + password + "@" + host)
	DB, err = sql.Open("mysql", username+":"+password+"@"+host)
	if err != nil {
		return err
	}
	return nil
}
