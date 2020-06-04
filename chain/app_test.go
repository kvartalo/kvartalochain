package chain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"kvartalochain/common"
	"kvartalochain/storage"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func setDbBalance(db *storage.Storage, addr common.Address, balance uint64) {
	var balanceBytes [8]byte
	binary.LittleEndian.PutUint64(balanceBytes[:], balance)
	db.Set(addr[:], balanceBytes[:])
}
func simulateTx(kApp *KvartaloABCI, sk *common.PrivateKey, from, to common.Address, amount, nonce uint64) (uint32, error) {
	// create and sign tx
	tx := common.NewTx(from, to, amount, nonce)
	sk.SignTx(tx)
	// addr := sk.Public().Address()
	if !common.VerifySignatureTx(tx) {
		return 4, fmt.Errorf("VerifySignatureTx failed")
	}
	txHex := hex.EncodeToString(tx.Bytes())

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
		balance := storage.GetBalance(kApp.db, addr)
		fmt.Println("	addr:", addr, " balance:", balance)
	}
}

func TestKvartaloApplication(t *testing.T) {
	tmpDir, err := ioutil.TempDir("./", "tmpTest")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	assert.Nil(t, err)

	archiveDb, err := badger.Open(badger.DefaultOptions(tmpDir).WithLogger(nil))
	require.Nil(t, err)
	defer archiveDb.Close()

	// initialize keys
	a := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
	b := "8h3u7NfgvUJsHJgKDUKwwVL1iZd3cwRtntpTfJ5Mefz2"
	sk0 := common.ImportKeyString(a)
	pk0 := sk0.Public()
	addr0 := pk0.Address()
	assert.Equal(t, "DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN", addr0.String())
	sk1 := common.ImportKeyString(b)
	pk1 := sk1.Public()
	addr1 := pk1.Address()
	assert.Equal(t, "HzeXxgjb589tVBs991jAyLUX7wreSZvrWnRxdGQS4co2", addr1.String())

	// initialize balances
	setDbBalance(db, addr0, 10)
	setDbBalance(db, addr1, 10)

	kApp := NewKvartaloApplication(db, archiveDb)
	printBalances(t, kApp, addr0, addr1)

	// get balance
	balance := storage.GetBalance(kApp.db, addr0)
	assert.Equal(t, uint64(10), balance)

	// addr0 send to addr1
	code, err := simulateTx(kApp, sk0, addr0, addr1, 10, 0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), code)
	balance = storage.GetBalance(kApp.db, addr0)
	assert.Equal(t, uint64(0), balance)
	balance = storage.GetBalance(kApp.db, addr1)
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
	balance = storage.GetBalance(kApp.db, addr0)
	assert.Equal(t, uint64(0), balance)
	balance = storage.GetBalance(kApp.db, addr1)
	assert.Equal(t, uint64(20), balance)
}
