package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"worker"
	"worker/pkg/config"
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

	server := worker.NewServer(conf)
	server.Run(context.Background())
}
