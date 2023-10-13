package lib

import (
	"context"
	"crypto/ecdsa"

	// "go/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type W3 struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	Address    *common.Address
}

func ConnectWeb3(provider string, privateKey string) (*W3, error) {
	// client, err := rpc.Dial(provider)
	client, err := ethclient.Dial(provider)
	if err != nil {
		return nil, err
	}
	privateKeyEcdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	publicKeyEcdsa := privateKeyEcdsa.Public()
	publicKeyEcdsaFin := publicKeyEcdsa.(*ecdsa.PublicKey)

	address := crypto.PubkeyToAddress(*publicKeyEcdsaFin)

	// fmt.Println(publicKeyEcdsa, address)

	return &W3{
		Client:     client,
		PrivateKey: privateKeyEcdsa,
		Address:    &address,
	}, nil
}

// func (w3 *W3) Close() {
// 	w3.client.Close()
// }

func (w3 *W3) GetGasPrice() (*big.Int, error) {
	return w3.Client.SuggestGasPrice(context.Background())
}

func (w3 *W3) GetNonce(address string) (uint64, error) {
	return w3.Client.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func (w3 *W3) NewTransaction() error {
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce: w3.GetNonce(*w3.Address),
		To:    w3.Address})
	return nil
}
