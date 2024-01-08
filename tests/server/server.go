package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"worker"
	"worker/pkg/config"
	"worker/pkg/connection"
	"worker/pkg/event"
	work "worker/pkg/worker"
)

func main() {
	// file := "/home/zhoujun/code/jun3/golang/socket/config/local.yml"
	file := "/home/ubuntu/code/golang/Worker/config/local.yml"
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

	server := worker.NewServer(conf, work.WithHandle(handler))
	server.Run(context.Background())
}

func handler(c *connection.Connection, e event.Event) {
	fmt.Println("connection", c.ID, "workID", c.WorkID, "topic", e.Topic, "data", e.Data)
}
