package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"worker/pkg/client"
	"worker/pkg/codec"
	"worker/pkg/config"
	"worker/pkg/connection"
)

func main() {
	content, err := os.ReadFile("/home/zhoujun/code/jun3/golang/socket/config/local.yml")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var conf config.Config
	if err := yaml.Unmarshal(content, &conf); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	addrress := fmt.Sprintf("%v:%v", conf.TCP.Host, conf.TCP.Port)
	ci, err := client.New(
		addrress,
		connection.WithCodec(codec.NewDESECB(conf.TCP.Secret)),
	)

	if err != nil {
		panic(err)
	}

	for {
		ci.Send("hello", "world")
		time.Sleep(time.Second * 5)
	}

}
