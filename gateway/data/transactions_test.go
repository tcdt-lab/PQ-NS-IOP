package data

import "testing"

func TestGatewayTransaction_InitStep(t *testing.T) {
	gt := GatewayTransaction{}
	gtu, vb := createMockVerifierAndUserGateway()
	idGt, idVb, err := gt.InitStep(gtu, vb)
	if err != nil {
		t.Errorf("Error in GatewayTransaction.InitStep: %v", err)
	}
	t.Log(idGt, idVb)
}

func createMockVerifierAndUserGateway() (GatewayUser, Verifier) {

	gtUser := GatewayUser{
		Id:           1,
		Salt:         "test",
		Password:     "test",
		PublicKeyDsa: "test",
		SecretKeyDsa: "test",
		PublicKeyKem: string("test"),
		SecretKeyKem: string("test"),
		Dsa_scheme:   "test",
		Kem_scheme:   "test",
	}
	bVerfier := Verifier{
		Id:        1,
		Ip:        "test",
		Port:      "test",
		PublicKey: "test",
	}
	return gtUser, bVerfier
}
