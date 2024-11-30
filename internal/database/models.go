package database

type WalletRecord struct {
	ID           uint64 `db:"id"`
	PhoneNumber  string `db:"phone_number"`
	PublicKey    string `db:"public_key"`
	Pin          string `db:"pin"`
	KeystorePath string `db:"keystore_path"`
}
