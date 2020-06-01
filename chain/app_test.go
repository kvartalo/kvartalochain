package chain

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"kvartalochain/common"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func setDbBalance(db *badger.DB, addr common.Address, balance uint64) {
	var balanceBytes [8]byte
	binary.BigEndian.PutUint64(balanceBytes[:], balance)
	txn := db.NewTransaction(true)
	if err := txn.Set(addr[:], balanceBytes[:]); err == badger.ErrTxnTooBig {
		_ = txn.Commit()
	}
	_ = txn.Commit()
}

func TestKvartaloApplication(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/test-kvartalochain/"))
	require.Nil(t, err)
	defer db.Close()

	// initialize keys
	a := "2jRarmCqNmQJw3gCfDZQEc2Yjhanq38VtW3Ua1toFSFbvQw1U1FUd6ojYcf6Lwbbr52qZMekyPKoXqVjk5a6enrP"
	b := "5BvqkWp8r8ZiLFk8DTkR2Ju1x6X19yrBE5Z8eGqKjXTRjTspqUbApkpZtBLRgYP2V7YLNWYwt2piur8DimjJMFCh"
	sk0 := common.ImportKeyString(a)
	pk0 := sk0.Public()
	addr0 := pk0.Address()
	assert.Equal(t, "GDB8nUdRW1dM4AhQVmnxoFdbTgvRs8YWSmNvwUgzK2K1", addr0.String())
	sk1 := common.ImportKeyString(b)
	pk1 := sk1.Public()
	addr1 := pk1.Address()
	assert.Equal(t, "9D42S4YDHbkdqL38gscQCHABj4nJYMTALKEFvZPfix6V", addr1.String())

	// initialize balances at 100
	setDbBalance(db, addr0, 100)
	setDbBalance(db, addr1, 100)

	kApp := NewKvartaloApplication(db)
	fmt.Println()

	// get balance
	balance, err := kApp.getBalance(addr0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(100), balance)

	// create and sign tx
	tx := common.NewTx(addr0, addr1, 10)
	sk0.SignTx(tx)
	assert.NotEqual(t, []byte{}, tx.Signature)
	assert.True(t, common.VerifySignatureTx(pk0, tx))
	txStr, err := json.Marshal(tx)
	assert.Nil(t, err)
	txHex := hex.EncodeToString(txStr)

	// validate Tx
	code := kApp.isValid([]byte(txHex))
	assert.Equal(t, uint32(0), code)

	// DeliverTx
	_ = kApp.BeginBlock(abcitypes.RequestBeginBlock{})
	req := abcitypes.RequestDeliverTx{
		Tx: []byte(txHex),
	}
	_ = kApp.DeliverTx(req)
	_ = kApp.Commit()
	balance, err = kApp.getBalance(addr0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(90), balance)
	balance, err = kApp.getBalance(addr1)
	assert.Nil(t, err)
	assert.Equal(t, uint64(110), balance)
}
