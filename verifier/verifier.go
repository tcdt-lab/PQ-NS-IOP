package main

import (
	"go.uber.org/zap"
	"os"
	config "verifier/config"
	"verifier/data_access"
	"verifier/logic"
	"verifier/network"
	"verifier/utility"
)

const PASSWORD = "password"

func main() {
	config, err := config.ReadYaml()
	db, err := utility.GetDBConnection(*config)
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
	adminId, bootstrapId, err := logic.InitStepLogic(db)
	if err != nil {
		zap.L().Error("Error while generating keys", zap.Error(err))
		os.Exit(1)
	}
	cacheHandler := data_access.NewCacheHandlerDA()
	cacheHandler.SetUserAdminId(adminId)
	cacheHandler.SetBootstrapVerifierId(bootstrapId)

	if err != nil {
		zap.L().Error("Error while getting db connection", zap.Error(err))
		os.Exit(1)
	}
	network.StartServer(config, db)
}
