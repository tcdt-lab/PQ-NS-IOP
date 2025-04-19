package main

import (
	"database/sql"
	"fmt"
	gtConfig "gateway/config"
	gtStateMachine "gateway/logic/state_machines"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"simulation/dsa"
	"simulation/key_distribution"
	"simulation/trust"
	vConfig "verifier/config"
	vStateMachine "verifier/logic/state_machines"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	fmt.Println("Enter the right input:")
	fmt.Println("1. Key Distribution Sequential")
	fmt.Println("2. Key Distribution Parallel")
	fmt.Println("3. DSA Sequential")
	fmt.Println("4. DSA Parallel")
	fmt.Println("5. Trust Score Calculation")
	fmt.Println("6. Key Distribution Operation(GV)")
	fmt.Println("7. Key Distribution Operation(VV)")
	fmt.Println("8. Get Info Operation (GV) ")
	fmt.Println("9. Get Info Operation (VV) ")
	fmt.Println("10. Balance Check Process")
	fmt.Println("Input:")
	var i int
	fmt.Scanf("%d", &i)
	runOperation(i)
	//trust.SimulateTrustScoreCalculation()
	//key_distribution.KeyDistibutionSequencial()
	//dsa.RunParallelDSA()
}

func runOperation(i int) {
	reqNum, _ := util.GenerateRequestNumber()

	switch i {
	case 1:
		key_distribution.KeyDistibutionSequencial()
	case 2:
		key_distribution.KeyDistributionParalell()
	case 3:
		dsa.RunSequentialDSA()
	case 4:
		dsa.RunParallelDSA()
	case 5:
		trust.SimulateTrustScoreCalculation()
	case 6:
		db, _ := getDbConnectionGT()
		boostrapKeyStateMachine := gtStateMachine.GenerateKEyDistroStateMachine(reqNum, db)
		boostrapKeyStateMachine.Transit()
	case 7:
		db, _ := getDbConnectionV()
		boostrapKeyStateMachine := vStateMachine.GenerateKeyDistroStateMachine(reqNum, db)
		boostrapKeyStateMachine.Transit()
	case 8:
		db, _ := getDbConnectionGT()
		boostrapGetInfoStateMachine := gtStateMachine.GenerateBootstrapGentInfoStateMachine(reqNum, db)
		boostrapGetInfoStateMachine.Transit()
	case 9:
		db, _ := getDbConnectionV()
		boostrapGetInfoStateMachine := vStateMachine.GenerateBootstrapGentInfoStateMachine(reqNum, db)
		boostrapGetInfoStateMachine.Transit()
	case 10:
		destinatinIp := "127.0.0.1"
		destinatinPort := "50052"
		db, _ := getDbConnectionGT()
		boostrapKeyStateMachine := gtStateMachine.GenerateBalanceCheckStateMachine(reqNum, destinatinIp, destinatinPort, db)
		boostrapKeyStateMachine.Transit()
	}
}

func getDbConnectionGT() (*sql.DB, error) {
	c, err := getConfigGT()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getDbConnectionV() (*sql.DB, error) {
	c, err := getConfigV()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getConfigV() (vConfig.Config, error) {
	cfg, err := vConfig.ReadYaml()
	return *cfg, err
}
func getConfigGT() (gtConfig.Config, error) {
	cfg, err := gtConfig.ReadYaml()
	return *cfg, err
}
