package trust

import (
	"math"
	"verifier/config"
)

type BaselineTrustScoreCalculator struct {
	c config.Config
}

func (b *BaselineTrustScoreCalculator) CalculateTrustScore(lastValidationResult float64, validationResultHistory []float64, scoresHistory []float64) (float64, float64, float64) {
	discountFactor := b.c.Trust.BaselineDiscountFactor
	if len(scoresHistory) == 0 {
		return ((1 - discountFactor) * lastValidationResult), discountFactor, 0
	}
	sum := float64(0)

	for index, vrh := range validationResultHistory {

		sum += vrh * math.Pow(discountFactor, float64(len(validationResultHistory)-(index+1)))
	}
	trstScore := (1 - discountFactor) * sum
	//trstScore := (scoresHistory[len(scoresHistory)-1] * discountFactor) + ((1 - discountFactor) * lastValidationResult)
	return trstScore, discountFactor, 0
}
