package asymmetric

import (
	ecc "cryptography_helper/pkg/asymmetric/ecc"
	pq "cryptography_helper/pkg/asymmetric/pq"
)

type AsymmetricHandler interface {
	DSKeyGen(schemeName string) (string, string, error)
	Sign(sKey string, message []byte, schemeName string) ([]byte, error)
	Verify(pkey string, message []byte, signature []byte, schemeName string) (bool, error)

	KEMKeyGen(schemeName string) (string, string, error)
	KemGenerateSecretKey(secKey string, pubKey string, ciphertext string, schemeName string) ([]byte, []byte, error)
}

func NewAsymmetricHandler(cryptographyType string) AsymmetricHandler {

	switch cryptographyType {
	case "ECC":
		return &ecc.Ecc_handler{&ecc.Ecc_dh{}, &ecc.Ecc_dsa{}}
	case "PQ":
		return &pq.Pq_handler{Mldsa: &pq.MLDSA{}, Mlkem: &pq.MLKEM{}}
	default:
		return nil
	}
}
