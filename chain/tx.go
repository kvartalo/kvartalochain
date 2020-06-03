package chain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"kvartalochain/common"
	"kvartalochain/storage"
)

const ERRFORMAT = uint32(1)
const ERRDB = uint32(2)
const ERRNONCE = uint32(3)
const ERRNOFUNDS = uint32(4)

func (app *KvartaloABCI) isValid(txRaw []byte) (code uint32) {
	txBytes, err := hex.DecodeString(string(txRaw))
	if err != nil {
		return ERRFORMAT // invalid tx format
	}

	tx, err := common.TxFromBytes(txBytes)
	if err != nil {
		fmt.Println("AI", err)
		return ERRFORMAT // invalid tx format
	}

	senderBalance, err := storage.GetBalance(app.db, tx.From)
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

func (app KvartaloABCI) performTx(txRaw []byte) uint32 {
	txBytes, err := hex.DecodeString(string(txRaw))
	if err != nil {
		return ERRFORMAT // invalid tx format
	}
	tx, err := common.TxFromBytes(txBytes)
	if err != nil {
		return ERRFORMAT // invalid tx format
	}

	// TODO check signature

	dbNonce, err := storage.GetNonce(app.db, tx.From)
	if err != nil {
		return ERRNONCE // error getting nonce
	}
	if dbNonce != tx.Nonce {
		return ERRNONCE
	}

	senderBalance, err := storage.GetBalance(app.db, tx.From)
	if err != nil {
		return ERRDB // error getting balance
	}
	receiverBalance, err := storage.GetBalance(app.db, tx.To)
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
	err = app.currentBatch.Set(append(storage.PREFIXNONCE, tx.From[:]...), newNonce[:])
	if err != nil {
		return ERRDB
	}

	// if node is in 'archive' mode, store history of tx
	if app.archive {
		// format in DB:
		// 	key: PREFIXHISTORY | address | count
		// 	value: tx.Bytes()

		// store tx for sender
		txFromCount, err := storage.GetTxCount(app.db, tx.From)
		if err != nil {
			return ERRDB
		}
		key := append(storage.PREFIXHISTORY, tx.From[:]...)
		var txFromCountBytes [8]byte
		binary.LittleEndian.PutUint64(txFromCountBytes[:], txFromCount)
		key = append(key, txFromCountBytes[:]...)
		err = app.currentBatch.Set(key, tx.Bytes())
		if err != nil {
			return ERRDB
		}
		txFromCountKey := append(storage.PREFIXHISTORY, tx.From[:]...)
		var txFromCountBytesNew [8]byte
		binary.LittleEndian.PutUint64(txFromCountBytesNew[:], txFromCount+1)
		err = app.currentBatch.Set(txFromCountKey, txFromCountBytesNew[:])
		if err != nil {
			return ERRDB
		}
		// store tx for receiver
		txToCount, err := storage.GetTxCount(app.db, tx.To)
		if err != nil {
			return ERRDB
		}
		key = append(storage.PREFIXHISTORY, tx.To[:]...)
		var txToCountBytes [8]byte
		binary.LittleEndian.PutUint64(txToCountBytes[:], txToCount)
		key = append(key, txToCountBytes[:]...)
		err = app.currentBatch.Set(key, tx.Bytes())
		if err != nil {
			return ERRDB
		}
		txToCountKey := append(storage.PREFIXHISTORY, tx.To[:]...)
		var txToCountBytesNew [8]byte
		binary.LittleEndian.PutUint64(txToCountBytesNew[:], txToCount+1)
		err = app.currentBatch.Set(txToCountKey, txToCountBytesNew[:])
		if err != nil {
			return ERRDB
		}
	}

	// fmt.Println("addr:", tx.From.String(), " balance: ", newSenderBalance)
	// fmt.Println("addr:", tx.To.String(), " balance: ", newReceiverBalance)
	return 0
}
