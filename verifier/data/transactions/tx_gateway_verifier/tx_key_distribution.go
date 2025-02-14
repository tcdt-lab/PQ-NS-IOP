package tx_gateway_verifier

import (
	"context"
	"database/sql"
	"verifier/data"
)

func SharedKeyAndGatewayRegistration(verifierUser data.VerifierUser, gateway data.Gateway, db *sql.DB) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO gateways (Ip, Port, Public_Key_Kem, Public_Key_Sig, Kem_Scheme, Sig_Scheme, ticket, Symmetric_Key) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", gateway.Ip, gateway.Port, gateway.PublicKeyKem, gateway.PublicKeySig, gateway.KemScheme, gateway.SigScheme, gateway.Ticket, gateway.SymmetricKey)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE verifier_user SET salt=?, password=?, secret_key_sig=?, public_key_sig=?, secret_key_kem=?, public_key_kem=?,symmetric_key=?  WHERE id=?", verifierUser.Salt, verifierUser.Password, verifierUser.SecretKeySig, verifierUser.PublicKeySig, verifierUser.SecretKeyKem, verifierUser.PublicKeyKem, verifierUser.SymmetricKey, verifierUser.Id)
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
