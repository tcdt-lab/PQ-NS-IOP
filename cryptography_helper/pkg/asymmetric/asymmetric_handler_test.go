package asymmetric

import "testing"

func TestNewAsymmetricHandler(t *testing.T) {
	util := NewAsymmetricHandler("ECC")
	t.Log(util)
}
