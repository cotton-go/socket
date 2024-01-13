package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/cotton-go/socket"
	"github.com/cotton-go/socket/pkg/config"
	"github.com/cotton-go/socket/pkg/connection"
	"github.com/cotton-go/socket/pkg/event"
	"github.com/cotton-go/socket/pkg/worker"
)

func main() {
	var file string
	flag.StringVar(&file, "conf", "", "配置文件")
	flag.Parse()
	if file == "" {
		fmt.Println("配置文件不能为空")
		os.Exit(0)
	}

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

	server := socket.NewServer(conf, worker.WithHandle(handler))
	server.Run(context.Background())
}

func handler(c *connection.Connection, e event.Event) {
	fmt.Println("connection", c.ID, "workID", c.WorkID, "topic", e.Topic, "data", e.Data)
	if e.Topic == event.TopicByLogin {
		c.Send("msg1", 1)
		c.Send("msg", map[string]string{"time": time.Now().String()})
		c.Send("msg", c)
	}
}
