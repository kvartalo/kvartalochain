package chain

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"kvartalochain/common"

	"github.com/dgraph-io/badger"
)

func (app *KvartaloABCI) isValid(txRaw []byte) (code uint32) {
	txBytes, err := hex.DecodeString(string(txRaw))
	if err != nil {
		fmt.Println("ASDF", string(txRaw))
		return 1 // invalid tx format
	}

	var tx common.Tx
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		return 1 // invalid tx format
	}

	fmt.Println("TXFROM", tx.From)
	senderBalance, err := app.getBalance(tx.From)
	if err != nil {
		return 2 // error getting balance
	}
	fmt.Println("balance", senderBalance)

	if senderBalance < tx.Amount {
		fmt.Println("[not enough funds] sender:", tx.From, "\nsenderBalance:", senderBalance, ", tx.Amount:", tx.Amount)
		return 3 // not enough funds
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
				balance = binary.BigEndian.Uint64(val)
				return err
			})
		}
		return nil
	})
	return balance, err
}

func (app KvartaloABCI) performTx(txRaw []byte) uint32 {
	txBytes, err := hex.DecodeString(string(txRaw))
	if err != nil {
		return 1 // invalid tx format
	}
	var tx common.Tx
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		return 1 // invalid tx format
	}

	// TODO check signature

	senderBalance, err := app.getBalance(tx.From)
	if err != nil {
		return 2 // error getting balance
	}
	receiverBalance, err := app.getBalance(tx.To)
	if err != nil {
		return 2 // error getting balance
	}

	// TODO add checks

	newSenderBalance := senderBalance - tx.Amount
	fmt.Println(receiverBalance)
	newReceiverBalance := receiverBalance + tx.Amount

	var newSenderBalanceBytes [8]byte
	binary.BigEndian.PutUint64(newSenderBalanceBytes[:], newSenderBalance)
	err = app.currentBatch.Set(tx.From[:], newSenderBalanceBytes[:])
	if err != nil {
		return 3
	}
	var newReceiverBalanceBytes [8]byte
	binary.BigEndian.PutUint64(newReceiverBalanceBytes[:], newReceiverBalance)
	err = app.currentBatch.Set(tx.To[:], newReceiverBalanceBytes[:])
	if err != nil {
		return 3
	}

	fmt.Println("addr:", tx.From.String(), " balance: ", newSenderBalance)
	fmt.Println("addr:", tx.To.String(), " balance: ", newReceiverBalance)
	return 0
}
