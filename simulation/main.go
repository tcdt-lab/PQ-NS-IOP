package main

import (
	"go.uber.org/zap"
	"simulation/key_distribution"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	//trust.SimulateTrustScoreCalculation()
	key_distribution.KeyDistibutionSequencial()
}
