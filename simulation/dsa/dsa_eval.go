package dsa

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

func RunSequentialDSA() {
	var elapsedTimes []float64
	var excelRows [][]float64
	c, _ := getConfig()
	db, _ := getDbConnection(c)
	for i := 0; i < 1000; i++ {
		println("*******************************************" + strconv.Itoa(i) + "*******************************************")
		println("start time: " + time.Now().String())
		elapseTime := runDSA(db)
		elapsedTimes = append(elapsedTimes, elapseTime)
		excelRows = append(excelRows, []float64{float64(i + 1), elapseTime})
	}
	excel_handler.WriteToAnExcelFile(c.Security.DSAScheme+"_DSA_elapsed_time_seq", excelRows)
	return
}

func RunParallelDSA() {
	var elapsedTimes []float64
	var excelRows [][]float64
	var avergeExcelRows [][]float64
	c, _ := getConfig()
	db, _ := getDbConnection(c)

	for j := 0; j < 9; j++ {
		var wg sync.WaitGroup
		var startSignal sync.WaitGroup
		startSignal.Add(1)
		var sum = float64(0)
		var array = [9]int{10, 20, 30, 40, 50, 60, 70, 80, 90}
		for i := 0; i < array[j]; i++ {
			wg.Add(1)

			go func(j int) {

				defer wg.Done()
				startSignal.Wait()
				println("start time: " + time.Now().String())
				println("*******************************************" + strconv.Itoa(i) + "*******************************************")
				elapseTime := runDSA(db)
				elapsedTimes = append(elapsedTimes, elapseTime)
				sum += elapseTime
				excelRows = append(excelRows, []float64{float64(j + 1), elapseTime})
			}(i)

		}

		startSignal.Done()
		wg.Wait()
		average := sum / float64(array[j])
		avergeExcelRows = append(avergeExcelRows, []float64{float64(array[j]), average})
		excel_handler.WriteToAnExcelFile(c.Security.DSAScheme+"_DSA_elapsed_time_parallel", excelRows)
	}
	excel_handler.WriteToAnExcelFile("avg"+c.Security.DSAScheme+"_DSA_elapsed_time_parallel", avergeExcelRows)
	return

}

func runDSA(db *sql.DB) float64 {
	startTime := time.Now()

	reqNum, err := util.GenerateRequestNumber()
	println("Request Number: " + strconv.FormatInt(reqNum, 10))
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}

	balanceCheckStateMachine := state_machines.GenerateBalanceCheckStateMachineForEvalDSA(reqNum, "127.0.0.1", "50090", db)
	balanceCheckStateMachine.Transit()
	elapsedTime := time.Since(startTime)
	return elapsedTime.Seconds()
}

func getConfig() (config.Config, error) {
	cfg, err := config.ReadYaml()
	return *cfg, err
}

func getDbConnection(c config.Config) (*sql.DB, error) {

	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}
