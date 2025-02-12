package main

import (
	"fmt"
	"gateway/config"
	"gateway/data_access"
	"gateway/logic"
	"gateway/logic/state_machines"
	"gateway/message_handler/util"
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

	//Init Steps
	//1.Generate DSA and KEM KEYS and Save Admin and Bootstrap Verifier data
	BootLogic(c)
	//2. phase one of the protocol
	PhaseOneExecute()

	//3. phase two of the protocol

	if err != nil {
		zap.L().Error("Error while generating keys", zap.Error(err))
		os.Exit(1)
	}

}

func BootLogic(c *config.Config) {
	bootstrapId, adminId, err := logic.InintStepLogic(c)
	if err != nil {
		fmt.Print(err)
		zap.L().Error("Error while generating keys", zap.Error(err))
		os.Exit(1)
	}
	zap.L().Info("BootstrapId and AdminId", zap.Int64("BootstrapId", bootstrapId), zap.Int64("AdminId", adminId))
	cacheHandler := data_access.NewCacheHandlerDA()
	cacheHandler.SetBootstrapVerifierId(int(bootstrapId))
	cacheHandler.SetUserAdminId(int(adminId))
}

func PhaseOneExecute() {
	reqNum, err := util.GenerateRequestNumber()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}

	boostrapKeyStateMachine := state_machines.GenerateBootstrapStateMachine(reqNum)
	boostrapKeyStateMachine.Transit()
}
