package common

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
)

var debug = true

func TestNewKey(t *testing.T) {
	pk, sk, err := NewKey()
	assert.Nil(t, err)
	if debug {
		fmt.Println("sk", sk.String())
		fmt.Println("pk", pk.String())
		fmt.Println("address", pk.Address().String())
	}
}

func TestImportPrivateKey(t *testing.T) {
	skStr := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
	sk := ImportKey(base58.Decode(skStr))
	pk := sk.Public()
	assert.Equal(t, "24vMSutZWU7KkraGG8LRzELJH2nvucgUpuDfE89pnyrYB", pk.String())
	assert.Equal(t, "DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN", pk.Address().String())

	skStr = "8h3u7NfgvUJsHJgKDUKwwVL1iZd3cwRtntpTfJ5Mefz2"
	sk = ImportKeyString(skStr)
	pk = sk.Public()
	assert.Equal(t, "xmvjz2h1VLHSsftHi4Bx6vxsjpp3xRZZyWPBELSgC237", pk.String())
	assert.Equal(t, "HzeXxgjb589tVBs991jAyLUX7wreSZvrWnRxdGQS4co2", pk.Address().String())
}

func TestSign(t *testing.T) {
	skStr := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
	sk := ImportKey(base58.Decode(skStr))
	pk := sk.Public()
	addr := pk.Address()
	m := []byte("helloworld")
	sig, err := sk.HashAndSign(m)
	assert.Nil(t, err)
	assert.True(t, VerifySignature(&addr, m, sig))
	assert.True(t, !VerifySignature(&addr, []byte("fake"), sig))
	sig[0] = 0
	assert.True(t, !VerifySignature(&addr, m, sig))
}

func TestSignTxVerifyTx(t *testing.T) {
	skStr0 := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
	sk0 := ImportKeyString(skStr0)
	pk0 := sk0.Public()
	addr0 := pk0.Address()

	tx := &Tx{
		From:      addr0,
		To:        addr0,
		Amount:    10,
		Signature: []byte{},
	}

	sk0.SignTx(tx)
	assert.NotEqual(t, []byte{}, tx.Signature)

	assert.True(t, VerifySignatureTx(&addr0, tx))
}

func TestTxMarshalers(t *testing.T) {
	skStr0 := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
	sk0 := ImportKeyString(skStr0)
	pk0 := sk0.Public()
	addr0 := pk0.Address()

	tx := &Tx{
		From:      addr0,
		To:        addr0,
		Amount:    10,
		Nonce:     0,
		Signature: []byte{},
	}
	sk0.SignTx(tx)

	txStr, err := json.Marshal(tx)
	assert.Nil(t, err)

	var txParsed Tx
	err = json.Unmarshal(txStr, &txParsed)
	assert.Nil(t, err)
	assert.Equal(t, tx, &txParsed)
}
func TestAddressFromString(t *testing.T) {
	addr0, err := AddressFromString("GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1")
	assert.Nil(t, err)
	assert.Equal(t, "GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1", addr0.String())

	_, err = AddressFromString("GDB8nUdRW1dM4AhQ")
	assert.NotNil(t, err)

}
func TestAddressMarshallers(t *testing.T) {
	addr, err := AddressFromString("GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1")
	assert.Nil(t, err)
	assert.Equal(t, "GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1", addr.String())

	addrStr, err := json.Marshal(&addr)
	assert.Nil(t, err)
	assert.Equal(t, `"GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1"`, string(addrStr))

	var addrParsed Address
	err = json.Unmarshal(addrStr, &addrParsed)
	assert.Nil(t, err)
	assert.Equal(t, addr, addrParsed)
}
