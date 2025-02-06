package main

import (
	"fmt"
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

	fmt.Println("Hello World")

	//Init Steps
	//1.Generate DSA and KEM KEYS and Save Admin and Bootstrap Verifier data
	bootstrapId, adminId, err := logic.InintStepLogic(c)
	zap.L().Info("BootstrapId and AdminId", zap.Int64("BootstrapId", bootstrapId), zap.Int64("AdminId", adminId))

	//2. Run Boostrap Key Distribution State Machine

	if err != nil {
		zap.L().Error("Error while generating keys", zap.Error(err))
		os.Exit(1)
	}

}
