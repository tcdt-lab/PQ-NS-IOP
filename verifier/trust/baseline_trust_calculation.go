package trust

import "verifier/config"

type BaselineTrustScoreCalculator struct {
	c config.Config
}

func (b *BaselineTrustScoreCalculator) CalculateTrustScore(lastValidationResult float64, validationResultHistory []float64, scoresHistory []float64) float64 {
	discountFactor := b.c.Trust.BaselineDiscountFactor
	if len(scoresHistory) == 0 {
		return ((1 - discountFactor) * lastValidationResult)
	}
	trstScore := (scoresHistory[len(scoresHistory)-1] * discountFactor) + ((1 - discountFactor) * lastValidationResult)
	return trstScore
}
