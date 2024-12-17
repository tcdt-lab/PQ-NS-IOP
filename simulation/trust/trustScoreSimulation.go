package trust

import (
	"fmt"
	"simulation/excel_handler"
	"strconv"
	"test.org/verifier/trust"
	config2 "verifier/config"
)

const epochs = 500

func generateProofResults() []float64 {
	verfierdProof := 1.0
	unverfiedProof := 0.0

	var proofResultsV1 []float64

	for i := 0; i < epochs; i++ {
		if 100 < i && i < 250 {
			if i%6 == 0 {
				proofResultsV1 = append(proofResultsV1, unverfiedProof)
				continue
			} else {
				proofResultsV1 = append(proofResultsV1, verfierdProof)
				continue
			}
		}

		//if 120 < i && i < 150 {
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

	//trustScoreHistory = append(trustScoreHistory, float64(proofResults[0]))
	var excelRows [][]float64
	for i := 0; i < epochs; i++ {
		newTrustScore := trust.CalculateTrustScore(proofResults[i], proofResults[:i+1], trustScoreHistory)
		trustScoreHistory = append(trustScoreHistory, newTrustScore)
		excelRows = append(excelRows, []float64{float64(i + 1), newTrustScore})
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
	epsilonStr := strconv.FormatFloat(config.Trust.PopoviciuEpsilon, 'f', -1, 32)
	excel_handler.WriteToAnExcelFile(config.Trust.ScoreScheme+epsilonStr, excelRows)
}
