package trust

import (
	"fmt"
	"math"
	"verifier/config"
)

type PopoviciuTrustScoreCalculator struct {
	c config.Config
}

func (b *PopoviciuTrustScoreCalculator) CalculateTrustScore(lastValidationResult float64, validationResultHistory []float64, scoresHistory []float64) (float64, float64, float64) {
	discountFactor, variance := b.calculateDiscountFactor(validationResultHistory)
	if len(scoresHistory) == 0 || variance == 0 {
		return (lastValidationResult), discountFactor, variance
	}
	sum := float64(0)

	for index, score := range validationResultHistory {

		sum += score * math.Pow(discountFactor, float64(len(validationResultHistory)-(index+1)))
	}
	trstScore := (1 - discountFactor) * sum
	fmt.Println("************")
	fmt.Println("variance", variance)
	fmt.Println("sum", sum)
	fmt.Println("dicount", discountFactor)
	fmt.Println("trust", trstScore)
	fmt.Println("************")

	//trstScore := (scoresHistory[len(scoresHistory)-1] * discountFactor) + ((1 - discountFactor) * lastValidationResult)
	//fmt.Println("discount factor", len(scoresHistory), discountFactor)
	return trstScore, discountFactor, variance
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
	if len(scores) > 400 {
		fmt.Println("min", minElement)
		fmt.Println("max", maxElement)
	}
	return maxElement - minElement
}

func (b *PopoviciuTrustScoreCalculator) calculateDiscountFactor(scores []float64) (float64, float64) {
	index := len(scores)
	//if len(scores) > 50 {
	//	scores = scores[len(scores)-50:]
	//}
	variance := b.calculateVariance(scores)
	ExponentAdjustment := b.c.Trust.ExponentAdjustment
	if scores[index-1] == 0 {
		ExponentAdjustment = 15
	} else {
		ExponentAdjustment = 0.09
	}

	//rangeValue := b.calculateTrustRange(scores)
	//mean := b.calculateAverage(scores)
	rangeValue := float64(1)
	maxVariance := (math.Pow(rangeValue, 2)) / 4
	//discountFactor := 1 - (variance / (maxVariance + ExponentAdjustment))
	discountFactor := math.Exp(-1 * ExponentAdjustment * (variance / maxVariance))
	//discountFactor := 1 / (1 + math.Exp(ExponentAdjustment*(-1*mean)))
	fmt.Println("discount factor", index, ":", discountFactor)
	fmt.Println("variance", index, ":", variance)
	fmt.Println("max variance", index, ":", maxVariance)
	//fmt.Println(variance)
	return discountFactor, variance
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
