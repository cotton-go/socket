package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"worker/pkg/client"
	"worker/pkg/codec"
	"worker/pkg/config"
	"worker/pkg/connection"
	"worker/pkg/event"
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
		connection.WithHandle(func(c *connection.Connection, e event.Event) {
			fmt.Println("print msg", "topic", e.Topic, "data", e.Data)
		}),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("startd", "ci", ci)
	select {}
}
