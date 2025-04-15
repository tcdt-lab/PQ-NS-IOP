package main

import (
	"database/sql"
	"fmt"
	"gateway/config"
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
	db, _ := getDbConnection()
	//Init Steps
	//1.Generate DSA and KEM KEYS and Save Admin and Bootstrap Verifier data
	BootLogic(c, db)
	//2. phase one of the protocol
	PhaseOneExecute(db)

	//3. phase two of the protocol

	if err != nil {
		zap.L().Error("Error while generating keys", zap.Error(err))
		os.Exit(1)
	}

}

func BootLogic(c *config.Config, db *sql.DB) {
	_, _, err := logic.InintStepLogic(c, db)
	if err != nil {
		fmt.Print(err)
		zap.L().Error("Error in Inint Logic Step", zap.Error(err))
		os.Exit(1)
	}
}

func PhaseOneExecute(db *sql.DB) {
	zap.L().Info("Phase One Execution Started")
	reqNum, err := util.GenerateRequestNumber()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}

	boostrapKeyStateMachine := state_machines.GenerateKEyDistroStateMachine(reqNum, db)
	boostrapKeyStateMachine.Transit()

	reqNum, err = util.GenerateRequestNumber()
	boostrapGetInfoStateMachine := state_machines.GenerateBootstrapGentInfoStateMachine(reqNum, db)
	boostrapGetInfoStateMachine.Transit()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}
	zap.L().Info("Get Info State Machine Completed")

}

func BalanceCheckStart(db *sql.DB) {
	zap.L().Info("Balance Check Execution Started")
	reqNum, err := util.GenerateRequestNumber()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}
	balanceCheckStateMachine := state_machines.GenerateBalanceCheckStateMachineForEvalDSA(reqNum, "127.0.0.1", "50090", db)
	balanceCheckStateMachine.Transit()
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}
	zap.L().Info("Balance Check State Machine Completed")

}
func getConfig() (config.Config, error) {
	cfg, err := config.ReadYaml()
	return *cfg, err
}

func getDbConnection() (*sql.DB, error) {
	c, err := getConfig()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}
