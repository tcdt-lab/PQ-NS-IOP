package data

import (
	"context"
	"gateway/config"
	"gateway/utility"
)

type GatewayTransaction struct {
}

func (gt *GatewayTransaction) InitStep(user GatewayUser, bootstrapVerifier Verifier) (int64, int64, error) {

	cfg, err := config.ReadYaml()
	if err != nil {
		return 0, 0, err
	}
	db, err := utility.GetDBConnection(*cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return 0, 0, err
	}
	resGateway, err := tx.ExecContext(ctx, "INSERT INTO gateway_user (salt, password, public_key_dsa, secret_key_dsa, public_key_kem, secret_key_kem, dsa_scheme, kem_scheme) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.Salt, user.Password, user.PublicKeyDsa, user.SecretKeyDsa, user.PublicKeyKem, user.SecretKeyKem, user.Dsa_scheme, user.Kem_scheme)
	if err != nil {
		return 0, 0, err
	}
	gatewayId, _ := resGateway.LastInsertId()
	resVerifier, err := tx.ExecContext(ctx, "INSERT INTO verifiers (Ip, Port, public_key, Symmetric_Key) VALUES (?, ?, ?, ?)", bootstrapVerifier.Ip, bootstrapVerifier.Port, bootstrapVerifier.PublicKey, bootstrapVerifier.SymmetricKey)
	bootstarpVerfiierId, _ := resVerifier.LastInsertId()
	tx.Commit()

	return gatewayId, bootstarpVerfiierId, nil
}
