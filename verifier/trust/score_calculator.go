package trust

import "verifier/config"

type ScoreCalculator interface {
	CalculateTrustScore(lastValidationResult float64, validationResultHistory []float64, scoresHistory []float64) float64
}

type TrustScoreCalculatorFactory struct {
	c config.Config
}

func (t *TrustScoreCalculatorFactory) CreateTrustScoreCalculator(config config.Config) ScoreCalculator {
	scoreScheme := config.Trust.ScoreScheme
	if scoreScheme == "popoviciu" {
		return &PopoviciuTrustScoreCalculator{c: config}
	}
	if scoreScheme == "baseline" {
		return &BaselineTrustScoreCalculator{c: config}
	} else {
		return nil
	}
}
