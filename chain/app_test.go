package chain

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kvartalochain/common"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func setDbBalance(db *badger.DB, addr common.Address, balance uint64) {
	var balanceBytes [8]byte
	binary.LittleEndian.PutUint64(balanceBytes[:], balance)
	txn := db.NewTransaction(true)
	if err := txn.Set(addr[:], balanceBytes[:]); err == badger.ErrTxnTooBig {
		_ = txn.Commit()
	}
	_ = txn.Commit()
}
func simulateTx(kApp *KvartaloABCI, sk common.PrivateKey, from, to common.Address, amount, nonce uint64) (uint32, error) {
	// create and sign tx
	tx := common.NewTx(from, to, amount, nonce)
	sk.SignTx(tx)
	if !common.VerifySignatureTx(sk.Public(), tx) {
		return 4, fmt.Errorf("VerifySignatureTx failed")
	}
	txStr, err := json.Marshal(tx)
	if err != nil {
		return 4, err
	}
	txHex := hex.EncodeToString(txStr)

	// DeliverTx
	_ = kApp.BeginBlock(abcitypes.RequestBeginBlock{})
	req := abcitypes.RequestDeliverTx{
		Tx: []byte(txHex),
	}
	res := kApp.DeliverTx(req)
	_ = kApp.Commit()
	return res.Code, nil
}
func printBalances(t *testing.T, kApp *KvartaloABCI, addrs ...common.Address) {
	fmt.Println("balances:")
	for _, addr := range addrs {
		balance, err := kApp.getBalance(addr)
		assert.Nil(t, err)
		fmt.Println("	addr:", addr, " balance:", balance)
	}
}

func TestKvartaloApplication(t *testing.T) {
	tmpDir, err := ioutil.TempDir("./", "tmpTest")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)
	db, err := badger.Open(badger.DefaultOptions(tmpDir).WithLogger(nil))
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

	// initialize balances
	setDbBalance(db, addr0, 10)
	setDbBalance(db, addr1, 10)

	kApp := NewKvartaloApplication(db)
	printBalances(t, kApp, addr0, addr1)

	// get balance
	balance, err := kApp.getBalance(addr0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(10), balance)

	// addr0 send to addr1
	code, err := simulateTx(kApp, sk0, addr0, addr1, 10, 0)
	assert.Equal(t, uint32(0), code)
	balance, err = kApp.getBalance(addr0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), balance)
	balance, err = kApp.getBalance(addr1)
	assert.Nil(t, err)
	assert.Equal(t, uint64(20), balance)
	printBalances(t, kApp, addr0, addr1)

	// addr0 send to addr1 without funds
	code, err = simulateTx(kApp, sk0, addr0, addr1, 10, 1)
	assert.Nil(t, err)
	assert.Equal(t, ERRNOFUNDS, code) // expect not enough funds

	// addr1 send to addr0
	code, err = simulateTx(kApp, sk1, addr1, addr0, 10, 0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), code)
	printBalances(t, kApp, addr0, addr1) // addr0 and addr1 both have 10

	// addr0 send to addr1 and gets error because of nonce already used
	code, err = simulateTx(kApp, sk0, addr0, addr1, 10, 0)
	assert.Nil(t, err)
	assert.Equal(t, ERRNONCE, code)
	// addr0 send to addr1 with nonce 1 should work
	code, err = simulateTx(kApp, sk0, addr0, addr1, 10, 1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), code)
	printBalances(t, kApp, addr0, addr1)
	balance, err = kApp.getBalance(addr0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), balance)
	balance, err = kApp.getBalance(addr1)
	assert.Nil(t, err)
	assert.Equal(t, uint64(20), balance)

}
