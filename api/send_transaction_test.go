// +build integration

package api

import (
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"github.com/MinterTeam/minter-go-sdk/wallet"
	"math/big"
	"testing"
)

func TestApi_Send(t *testing.T) {
	data := transaction.NewSendData().
		SetCoin("MNT").
		SetValue(big.NewInt(0).Mul(big.NewInt(1), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))).
		MustSetTo("Mxee81347211c72524338f9680072af90744333146")

	newTransaction, err := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(data)
	if err != nil {
		t.Fatal(err)
	}

	mnemonic := "perfect bid satoshi giant cigar judge tonight possible harbor render else food"
	seed, err := wallet.Seed(mnemonic)
	if err != nil {
		t.Fatal(err)
	}
	mntWallet, err := wallet.NewWallet(seed)
	if err != nil {
		t.Fatal(err)
	}
	wantAddress := "Mxeeee1973381ab793719fff497b9a516719fcd5a2"
	if wantAddress != mntWallet.Address() {
		t.Fatalf("wallet address get %s, want %s", mntWallet.Address(), wantAddress)
	}

	wantedPrKey := "ecc067573863f893f2195d550ff0d703d31e5a1255791e410ca7ff7cf5f0a7aa"
	if mntWallet.PrivateKey() != wantedPrKey {
		t.Fatalf("privateKey get %s, want %s", mntWallet.PrivateKey(), wantedPrKey)
	}

	nonce, err := testApi.Nonce("Mxeeee1973381ab793719fff497b9a516719fcd5a2")
	if err != nil {
		t.Fatal(err)
	}

	signedTransaction, err := newTransaction.
		SetNonce(nonce).
		SetGasCoin("MNT").
		SetGasPrice(1).
		Sign(mntWallet.PrivateKey())
	if err != nil {
		t.Fatal(err)
	}

	res, err := testApi.SendTransaction(signedTransaction)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}
