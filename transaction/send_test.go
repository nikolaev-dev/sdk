package transaction

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestTransactionSend_Sign(t *testing.T) {
	value := big.NewInt(0).Mul(big.NewInt(1), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	address := "Mx1b685a7c1e78726c48f619c497a07ed75fe00483"
	symbolMNT := "MNT"
	data, err := NewSendData().
		SetCoin(symbolMNT).
		SetValue(value).
		SetTo(address)
	if err != nil {
		t.Fatal(err)
	}

	if string(data.Coin[:3]) != symbolMNT {
		t.Errorf("SendData.Coin got %s, want %s", data.Coin, symbolMNT)
	}

	addressBytes, err := hex.DecodeString(address[2:])
	if string(data.To[:]) != string(addressBytes) {
		t.Errorf("SendData.To got %s, want %s", string(data.To[:]), string(addressBytes))
	}

	if data.Value.String() != value.String() {
		t.Errorf("SendData.Value got %s, want %s", data.Value.String(), value.String())
	}
	tx, err := NewBuilder(TestNetChainID).NewTransaction(data)
	if err != nil {
		t.Fatal(err)
	}

	nonce := uint64(1)
	gasPrice := uint8(1)

	tx.SetNonce(nonce).SetGasPrice(gasPrice).SetGasCoin(symbolMNT)
	transaction := tx.(*object)

	if transaction.Nonce != nonce {
		t.Errorf("Nonce got %d, want %d", transaction.Nonce, nonce)
	}

	if transaction.ChainID != TestNetChainID {
		t.Errorf("ChainID got %d, want %d", transaction.ChainID, TestNetChainID)
	}

	if transaction.GasPrice != gasPrice {
		t.Errorf("GasPrice got %d, want %d", transaction.GasPrice, gasPrice)
	}

	gasCoinBytes := Coin{'\x4d', '\x4e', '\x54'} // MNT
	if string(transaction.GasCoin[:]) != string(gasCoinBytes[:]) {
		t.Errorf("GasCoin got %s, want %s", transaction.GasCoin, gasCoinBytes)
	}

	signedTx, err := transaction.Sign("07bc17abdcee8b971bb8723e36fe9d2523306d5ab2d683631693238e0f9df142")
	if err != nil {
		t.Fatal(err)
	}

	validSignature := "0xf8840102018a4d4e540000000000000001aae98a4d4e5400000000000000941b685a7c1e78726c48f619c497a07ed75fe00483880de0b6b3a7640000808001b845f8431ca01f36e51600baa1d89d2bee64def9ac5d88c518cdefe45e3de66a3cf9fe410de4a01bc2228dc419a97ded0efe6848de906fbe6c659092167ef0e7dcb8d15024123a"
	bytes, err := signedTx.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != validSignature {
		t.Errorf("EncodeTx got %s, want %s", string(bytes), validSignature)
	}
}
