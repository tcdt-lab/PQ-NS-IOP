package key_distribution

import (
	"database/sql"
	"gateway/config"
	"gateway/logic/state_machines"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"os"
	"simulation/excel_handler"
	"strconv"
	"sync"
	"time"
)

func KeyDistibutionSequencial() {
	var elapsedTimes []float64
	var excelRows [][]float64
	db, _ := getDbConnection()
	for i := 0; i < 1000; i++ {
		println("*******************************************" + strconv.Itoa(i) + "*******************************************")
		println("start time: " + time.Now().String())
		elapseTime := keyDistributionPQ(db)
		elapsedTimes = append(elapsedTimes, elapseTime)
		excelRows = append(excelRows, []float64{float64(i + 1), elapseTime})
	}
	excel_handler.WriteToAnExcelFile("pq_elapsed_time_seq", excelRows)
	return

}
func KeyDistributionParalell() {
	var elapsedTimes []float64
	var excelRows [][]float64
	db, _ := getDbConnection()
	var wg sync.WaitGroup
	var startSignal sync.WaitGroup
	startSignal.Add(1)
	for i := 0; i < 90; i++ {
		wg.Add(1)

		go func(j int) {

			defer wg.Done()
			startSignal.Wait()
			println("start time: " + time.Now().String())
			println("*******************************************" + strconv.Itoa(i) + "*******************************************")
			elapseTime := keyDistributionPQ(db)
			elapsedTimes = append(elapsedTimes, elapseTime)
			excelRows = append(excelRows, []float64{float64(j + 1), elapseTime})
		}(i)

	}
	startSignal.Done()
	wg.Wait()
	excel_handler.WriteToAnExcelFile("pq_elapsed_time", excelRows)
	return
}

func keyDistributionPQ(db *sql.DB) float64 {
	startTime := time.Now()

	reqNum, err := util.GenerateRequestNumber()
	println("Request Number: " + strconv.FormatInt(reqNum, 10))
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}

	boostrapKeyStateMachine := state_machines.GenerateKEyDistroStateMachine(reqNum, db)
	boostrapKeyStateMachine.Transit()
	elapsedTime := time.Since(startTime)
	return elapsedTime.Seconds()
}

func getConfig() (config.Config, error) {
	cfg, err := config.ReadYaml()
	return *cfg, err
}

func getDbConnection() (*sql.DB, error) {
	c, err := getConfig()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}
