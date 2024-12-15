package trust

import (
	"go.uber.org/zap"
	"verifier/config"
)

func CalculateTrustScore(lastValidationResult float64, trustScoreHistory []float64) float64 {
	// Calculate trust score
	c, err := config.ReadYaml()
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return 0
	}
	trustScoreCalculatorFactory := TrustScoreCalculatorFactory{}
	scoreCalculator := trustScoreCalculatorFactory.CreateTrustScoreCalculator(*c)
	return scoreCalculator.CalculateTrustScore(lastValidationResult, trustScoreHistory)
}
