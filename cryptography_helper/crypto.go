package main

import (
	"cryptography_helper/pkg/asymmetric"
	"cryptography_helper/pkg/asymmetric/pq"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	zap.L().Info("replaced zap's global loggers")
	var mldsa pq.MLDSA
	mldsa.KeyGen("ML-DSA-65")
	util := asymmetric.NewAsymmetricHandler("ECC")
	fmt.Println(util)
}
