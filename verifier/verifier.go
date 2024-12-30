package main

import (
	"go.uber.org/zap"
	"os"
	config "verifier/config"
	"verifier/network"
)

const PASSWORD = "password"

func main() {
	config, err := config.ReadYaml()
	if err != nil {
		panic(err)
	}

	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		os.Exit(1)
	}
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	network.StartServer(config)
}
