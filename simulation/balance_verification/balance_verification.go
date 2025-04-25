package balance_verification

import (
	"database/sql"
	"sync"

	"gateway/config"
	gtStateMachine "gateway/logic/state_machines"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"os"
	"simulation/excel_handler"
	"strconv"
	"time"
)

func RunSequencialBalanceVerification(destinationIp string, destinationPort string) {
	// This function is not implemented yet.
	// 1. Generate the proof and public inputs using the circuit
	// 2. Create a balance verification request
	// 3. Send the request to the gateway
	// 4. Wait for the response from the gateway
	// 5. Verify the proof using the verification key
	// 6. Return the result
	var elapsedTimes []float64
	var excelRows [][]float64
	c, _ := getConfig()
	db, _ := getDbConnection(c)
	for i := 0; i < 10; i++ {
		println("*******************************************" + strconv.Itoa(i) + "*******************************************")
		println("start time: " + time.Now().String())
		var mu sync.Mutex
		elapseTime := runBalanceVerification(db, destinationIp, destinationPort, &mu)
		elapsedTimes = append(elapsedTimes, elapseTime)
		excelRows = append(excelRows, []float64{float64(i + 1), elapseTime})
	}
	excel_handler.WriteToAnExcelFile(c.Security.DSAScheme+"_BV_elapsed_time_seq", excelRows)
	return
}

func RunParallelBalanceVerification(destinationIp string, destinationPort string) {
	var elapsedTimes []float64
	var excelRows [][]float64
	var avergeExcelRows [][]float64
	c, _ := getConfig()
	db, _ := getDbConnection(c)

	var mutex sync.Mutex
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
				elapseTime := runBalanceVerification(db, destinationIp, destinationPort, &mutex)
				elapsedTimes = append(elapsedTimes, elapseTime)
				sum += elapseTime
				excelRows = append(excelRows, []float64{float64(j + 1), elapseTime})
			}(i)

		}

		startSignal.Done()
		wg.Wait()
		average := sum / float64(array[j])
		avergeExcelRows = append(avergeExcelRows, []float64{float64(array[j]), average})
		excel_handler.WriteToAnExcelFile(c.Security.DSAScheme+"_BV_elapsed_time_parallel", excelRows)
	}
	excel_handler.WriteToAnExcelFile("avg"+c.Security.DSAScheme+"_BV_elapsed_time_parallel", avergeExcelRows)
	return

}
func runBalanceVerification(db *sql.DB, destinationIp string, destinationPort string, mutex *sync.Mutex) float64 {
	startTime := time.Now()

	reqNum, err := util.GenerateRequestNumber()
	println("Request Number: " + strconv.FormatInt(reqNum, 10))
	if err != nil {
		zap.L().Error("Error while generating request number", zap.Error(err))
		os.Exit(1)
	}

	boostrapKeyStateMachine := gtStateMachine.GenerateBalanceCheckStateMachine(reqNum, destinationIp, destinationPort, db, mutex)
	boostrapKeyStateMachine.Transit()
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
