package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/cotton-go/socket/pkg/client"
	"github.com/cotton-go/socket/pkg/codec"
	"github.com/cotton-go/socket/pkg/config"
	"github.com/cotton-go/socket/pkg/connection"
	"github.com/cotton-go/socket/pkg/event"
)

func main() {
	file := "/home/zhoujun/code/jun3/golang/socket/config/local.yml"
	// file := "/home/ubuntu/code/golang/Worker/config/local.yml"
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var conf config.Config
	if err := yaml.Unmarshal(content, &conf); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	ctx, cannel := context.WithCancel(context.Background())
	addrress := fmt.Sprintf("%v:%v", conf.TCP.Host, conf.TCP.Port)
	ci, err := client.New(
		addrress,
		connection.WithHandle(handler),
		connection.WithCodec(codec.NewDESECB(conf.TCP.Secret)),
		connection.WithClose(func(c *connection.Connection, e event.Event) {
			cannel()
		}),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("startd", "ci", ci)
	select {
	case <-ctx.Done():
		return
	}
}

func handler(c *connection.Connection, e event.Event) {
	fmt.Println("print msg", "topic", e.Topic, "data", e.Data)
}
