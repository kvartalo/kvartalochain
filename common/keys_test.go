package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
)

var debug = true

func TestNewKey(t *testing.T) {
	pk, sk, err := NewKey(rand.Reader)
	assert.Nil(t, err)
	if debug {
		fmt.Println("sk", base58.Encode(sk))
		fmt.Println("pk", pk.String())
		fmt.Println("address", pk.Address().String())
	}
}

func TestImportPrivateKey(t *testing.T) {
	skStr := "2jRarmCqNmQJw3gCfDZQEc2Yjhanq38VtW3Ua1toFSFbvQw1U1FUd6ojYcf6Lwbbr52qZMekyPKoXqVjk5a6enrP"
	sk := ImportKey(base58.Decode(skStr))
	pk := sk.Public()
	assert.Equal(t, "8f5pHm4Yzutru4UhoUPWxqnjn2PUJXHEDMngoPd8AoeB", pk.String())
	assert.Equal(t, "GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1", pk.Address().String())

	skStr = "5BvqkWp8r8ZiLFk8DTkR2Ju1x6X19yrBE5Z8eGqKjXTRjTspqUbApkpZtBLRgYP2V7YLNWYwt2piur8DimjJMFCh"
	sk = ImportKeyString(skStr)
	pk = sk.Public()
	assert.Equal(t, "3TSQkkbnuZQ4gPy8wass7Dp1VhyMwSzbjbZbGFAN9Xcy", pk.String())
	assert.Equal(t, "9D42S4YDHbkdqL38gscQCHABj4nJYMTALKEFvZPfix6V", pk.Address().String())
}

func TestSign(t *testing.T) {
	skStr := "2jRarmCqNmQJw3gCfDZQEc2Yjhanq38VtW3Ua1toFSFbvQw1U1FUd6ojYcf6Lwbbr52qZMekyPKoXqVjk5a6enrP"
	sk := ImportKey(base58.Decode(skStr))
	pk := sk.Public()
	m := []byte("helloworld")
	sig := sk.Sign(m)
	assert.True(t, VerifySignature(pk, m, sig))
	assert.True(t, !VerifySignature(pk, []byte("fake"), sig))
	sig[0] = 0
	assert.True(t, !VerifySignature(pk, m, sig))
}

func TestSignTxVerifyTx(t *testing.T) {
	skStr0 := "2jRarmCqNmQJw3gCfDZQEc2Yjhanq38VtW3Ua1toFSFbvQw1U1FUd6ojYcf6Lwbbr52qZMekyPKoXqVjk5a6enrP"
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

	assert.True(t, VerifySignatureTx(pk0, tx))
}

func TestTxMarshalers(t *testing.T) {
	skStr0 := "2jRarmCqNmQJw3gCfDZQEc2Yjhanq38VtW3Ua1toFSFbvQw1U1FUd6ojYcf6Lwbbr52qZMekyPKoXqVjk5a6enrP"
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
