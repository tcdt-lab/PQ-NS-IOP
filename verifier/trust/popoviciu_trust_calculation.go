package trust

import (
	"math"
	"verifier/config"
)

type PopoviciuTrustScoreCalculator struct {
	c config.Config
}

func (b *PopoviciuTrustScoreCalculator) CalculateTrustScore(lastValidationResult float64, scoresHistory []float64) float64 {
	discountFactor := b.calculateDiscountFactor(scoresHistory)
	trstScore := (scoresHistory[len(scoresHistory)-1] * discountFactor) + ((1 - discountFactor) * lastValidationResult)
	return trstScore
}

func (b *PopoviciuTrustScoreCalculator) calculateTrustRange(scores []float64) float64 {
	minElement := scores[0]
	maxElement := scores[0]
	for _, score := range scores {
		if score < minElement {
			minElement = score
		}
		if score > maxElement {
			maxElement = score
		}
	}
	return maxElement - minElement
}

func (b *PopoviciuTrustScoreCalculator) calculateDiscountFactor(scores []float64) float64 {
	variance := b.calculateVariance(scores)
	epsilon := b.c.Trust.PopoviciuEpsilon
	rangeValue := b.calculateTrustRange(scores)

	discountFactor := (variance/(math.Pow(rangeValue, 2)/4) + epsilon)
	return discountFactor
}

func (b *PopoviciuTrustScoreCalculator) calculateVariance(scores []float64) float64 {
	average := b.calculateAverage(scores)
	var sum float64
	for _, score := range scores {
		sum += math.Pow(score-average, 2)
	}
	return sum / float64(len(scores))
}

func (b *PopoviciuTrustScoreCalculator) calculateAverage(scores []float64) float64 {
	var sum float64
	for _, score := range scores {
		sum += score
	}
	return sum / float64(len(scores))
}
