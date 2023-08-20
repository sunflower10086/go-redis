package main

import (
	"fmt"
	"os"

	"github.com/sunflower10086/go-redis/config"
	"github.com/sunflower10086/go-redis/lib/logger"
	"github.com/sunflower10086/go-redis/resp/handler"
	"github.com/sunflower10086/go-redis/tcp"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "go-redis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, handler.NewRespHandler())

	if err != nil {
		panic(err)
	}
}
