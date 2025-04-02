package main

import (
	"database/sql"
	"go.uber.org/zap"
	"os"
	config "verifier/config"
	"verifier/data_access"
	"verifier/logic"
	"verifier/logic/state_machines"
	"verifier/message_handler/util"
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

	if isBootstrapNodeAvailable(config) {
		PhaseOneHandler(err, db)
	}
	network.StartServer(config, db)
}

func isBootstrapNodeAvailable(c *config.Config) bool {
	if c.BootstrapNode.Ip != "none" && c.BootstrapNode.Port != "none" {
		return true
	}
	return false
}

func PhaseOneHandler(err error, db *sql.DB) {
	reqId, err := util.GenerateRequestNumber()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}
	err = startKeyDistribution(db, reqId)
	if err != nil {
		zap.L().Error("Error while starting key distribution", zap.Error(err))
		os.Exit(1)
	}
	reqId, err = util.GenerateRequestNumber()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)

	}
	err = startGetInfo(db, reqId)
	if err != nil {
		zap.L().Error("Error while starting get info", zap.Error(err))
		os.Exit(1)
	}
	zap.L().Info("Phase One Execution Completed")
}

func startKeyDistribution(db *sql.DB, reqId int64) error {

	keyDistroFSM := state_machines.GenerateKeyDistroStateMachine(reqId, db)
	err := keyDistroFSM.Transit()
	if err != nil {
		zap.L().Error("Error while transiting key distribution state machine", zap.Error(err))
		return err
	}
	return nil
}

func startGetInfo(db *sql.DB, reqId int64) error {
	getInfoFSM := state_machines.GenerateBootstrapGentInfoStateMachine(reqId, db)
	err := getInfoFSM.Transit()
	if err != nil {
		zap.L().Error("Error while transiting get info state machine", zap.Error(err))
		return err
	}
	return nil
}
