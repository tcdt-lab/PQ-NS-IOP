package trust

import "verifier/config"

type BaselineTrustScoreCalculator struct {
	c config.Config
}

func (b *BaselineTrustScoreCalculator) CalculateTrustScore(lastValidationResult float64, scoresHistory []float64) float64 {
	discountFactor := b.c.Trust.BaselineDiscountFactor
	trstScore := (scoresHistory[len(scoresHistory)-1] * discountFactor) + ((1 - discountFactor) * lastValidationResult)
	return trstScore
}
