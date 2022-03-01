package master

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var G_config *Config

type Config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	MongodbUri      string   `json:"mongodbUri"`
}

func IntiConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	conf := Config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return err
	}
	G_config = &conf
	log.Println("配置初始化成功。。。")
	return nil
}
