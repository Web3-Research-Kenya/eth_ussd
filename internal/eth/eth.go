package eth

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	ethc "github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func Connect() *ethc.Client {
	client, err := ethc.Dial("/tmp/geth.ipc")
	if err != nil {
		log.Fatalf("failed to dial geth: %e\n", err)
	}
	return client
}

func NewWallet() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate private key: %e\n", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("cannot asert type: %e\n", err)
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil))[12:])

}

func NewKeyStore(pin string) (*accounts.Account, error) {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	password := pin

	account, err := ks.NewAccount(password)
	if err != nil {
		return nil, err
	}
	return &account, nil

}

func ImportKeystore(filePath, pin string) (*accounts.Account, error) {
	ks := keystore.NewKeyStore("/tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	account, err := ks.Import(jsonBytes, pin, pin)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Remove("/tmp/*"); err != nil {
		log.Fatal(err)
	}

	return &account, nil
}
