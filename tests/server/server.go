package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/cotton-go/socket"
	"github.com/cotton-go/socket/pkg/config"
	"github.com/cotton-go/socket/pkg/connection"
	"github.com/cotton-go/socket/pkg/event"
	"github.com/cotton-go/socket/pkg/worker"
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

	server := socket.NewServer(conf, worker.WithHandle(handler))
	server.Run(context.Background())
}

func handler(c *connection.Connection, e event.Event) {
	fmt.Println("connection", c.ID, "workID", c.WorkID, "topic", e.Topic, "data", e.Data)
}
