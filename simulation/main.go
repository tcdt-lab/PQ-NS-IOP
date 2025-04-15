package main

import (
	"go.uber.org/zap"
	"simulation/dsa"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	//trust.SimulateTrustScoreCalculation()
	//key_distribution.KeyDistibutionSequencial()
	dsa.RunParallelDSA()
}
