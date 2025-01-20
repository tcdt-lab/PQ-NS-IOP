package main

import (
	"gateway/config"
	"gateway/logic"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
)

func main() {
	c, err := config.ReadYaml()
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		os.Exit(1)
	}
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	//	KeyGen
	bootstrapId, adminId, err := logic.KeyGenLogic(c)
	zap.L().Info("BootstrapId and AdminId", zap.Int64("BootstrapId", bootstrapId), zap.Int64("AdminId", adminId))

	//keyDistribution

	if err != nil {
		zap.L().Error("Error while generating keys", zap.Error(err))
		os.Exit(1)
	}

}
