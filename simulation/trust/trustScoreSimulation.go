package trust

import (
	"fmt"
	"simulation/excel_handler"
	"strconv"
	"test.org/verifier/trust"
	config2 "verifier/config"
)

const epochs = 200

func generateProofResults() []float64 {
	verfierdProof := 1.0
	unverfiedProof := 0.0

	var proofResultsV1 []float64

	for i := 0; i < epochs; i++ {

		if 39 < i && i < 80 {
			if i%4 == 0 {
				proofResultsV1 = append(proofResultsV1, unverfiedProof)
				continue
			} else {
				proofResultsV1 = append(proofResultsV1, verfierdProof)
				continue
			}

			//proofResultsV1 = append(proofResultsV1, unverfiedProof)
			//continue
		}

		//if 81 < i && i < 111 {
		//if i%6 == 0 {
		//	proofResultsV1 = append(proofResultsV1, unverfiedProof)
		//	continue
		//} else {
		//	proofResultsV1 = append(proofResultsV1, verfierdProof)
		//	continue
		//}

		//	proofResultsV1 = append(proofResultsV1, unverfiedProof)
		//	continue
		//}
		proofResultsV1 = append(proofResultsV1, verfierdProof)
	}

	return proofResultsV1

}

func SimulateTrustScoreCalculation() {
	proofResults := generateProofResults()

	var trustScoreHistory []float64
	var discountFactors []float64
	var variances []float64

	//trustScoreHistory = append(trustScoreHistory, float64(proofResults[0]))
	var excelRows [][]float64
	for i := 0; i < epochs; i++ {
		newTrustScore, discountFactor, variance := trust.CalculateTrustScore(proofResults[i], proofResults[:i+1], trustScoreHistory)
		trustScoreHistory = append(trustScoreHistory, newTrustScore)
		discountFactors = append(discountFactors, discountFactor)
		variances = append(variances, variance)
		excelRows = append(excelRows, []float64{float64(i + 1), newTrustScore, discountFactor, variance})
	}
	for index, value := range trustScoreHistory {
		fmt.Printf("Index: %d, Value: %.20f\n", index, value)
	}
	config, err := config2.ReadYaml()
	if err != nil {
		fmt.Println(err)
		return
	}
	if config.Trust.ScoreScheme == "baseline" {
		discountStr := strconv.FormatFloat(config.Trust.BaselineDiscountFactor, 'f', -1, 32)
		excel_handler.WriteToAnExcelFile(config.Trust.ScoreScheme+discountStr, excelRows)
		return
	}
	exponentAdjustmentStr := strconv.FormatFloat(config.Trust.ExponentAdjustment, 'f', -1, 32)
	excel_handler.WriteToAnExcelFile(config.Trust.ScoreScheme+exponentAdjustmentStr, excelRows)
}
