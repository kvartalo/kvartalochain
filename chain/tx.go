package chain

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"kvartalochain/common"

	"github.com/dgraph-io/badger"
)

var NONCEKEY = []byte("nonce")

const ERRFORMAT = uint32(1)
const ERRDB = uint32(2)
const ERRNONCE = uint32(3)
const ERRNOFUNDS = uint32(4)

func (app *KvartaloABCI) isValid(txRaw []byte) (code uint32) {
	txBytes, err := hex.DecodeString(string(txRaw))
	if err != nil {
		return ERRFORMAT // invalid tx format
	}

	var tx common.Tx
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		return ERRFORMAT // invalid tx format
	}

	senderBalance, err := app.getBalance(tx.From)
	if err != nil {
		return ERRDB // error getting balance
	}

	if senderBalance < tx.Amount {
		fmt.Println("[not enough funds] sender:", tx.From, "\nsenderBalance:", senderBalance, ", tx.Amount:", tx.Amount)
		return ERRNOFUNDS // not enough funds
	}

	// return 0 code if valid
	return code
}

func (app KvartaloABCI) getBalance(addr common.Address) (uint64, error) {
	var balance uint64
	// TODO when addr not found in db, return balance 0
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(addr[:])
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return item.Value(func(val []byte) error {
				balance = binary.LittleEndian.Uint64(val)
				return err
			})
		}
		return nil
	})
	return balance, err
}

func (app KvartaloABCI) getNonce(addr common.Address) (uint64, error) {
	var nonce uint64
	err := app.db.View(func(txn *badger.Txn) error {
		nonceKey := append(NONCEKEY, addr[:]...)
		item, err := txn.Get(nonceKey)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return item.Value(func(val []byte) error {
				nonce = binary.LittleEndian.Uint64(val)
				return err
			})
		}
		return nil
	})
	return nonce, err
}

func (app KvartaloABCI) performTx(txRaw []byte) uint32 {
	txBytes, err := hex.DecodeString(string(txRaw))
	if err != nil {
		return ERRFORMAT // invalid tx format
	}
	var tx common.Tx
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		return ERRFORMAT // invalid tx format
	}

	// TODO check signature

	dbNonce, err := app.getNonce(tx.From)
	if err != nil {
		return ERRNONCE // error getting nonce
	}
	if dbNonce != tx.Nonce {
		return ERRNONCE
	}

	senderBalance, err := app.getBalance(tx.From)
	if err != nil {
		return ERRDB // error getting balance
	}
	receiverBalance, err := app.getBalance(tx.To)
	if err != nil {
		return ERRDB // error getting balance
	}

	// TODO add checks

	newSenderBalance := senderBalance - tx.Amount
	newReceiverBalance := receiverBalance + tx.Amount

	var newSenderBalanceBytes [8]byte
	binary.LittleEndian.PutUint64(newSenderBalanceBytes[:], newSenderBalance)
	err = app.currentBatch.Set(tx.From[:], newSenderBalanceBytes[:])
	if err != nil {
		return ERRDB
	}
	var newReceiverBalanceBytes [8]byte
	binary.LittleEndian.PutUint64(newReceiverBalanceBytes[:], newReceiverBalance)
	err = app.currentBatch.Set(tx.To[:], newReceiverBalanceBytes[:])
	if err != nil {
		return ERRDB
	}

	var newNonce [8]byte
	binary.LittleEndian.PutUint64(newNonce[:], dbNonce+1)
	err = app.currentBatch.Set(append(NONCEKEY, tx.From[:]...), newNonce[:])
	if err != nil {
		return ERRDB
	}

	// fmt.Println("addr:", tx.From.String(), " balance: ", newSenderBalance)
	// fmt.Println("addr:", tx.To.String(), " balance: ", newReceiverBalance)
	return 0
}
