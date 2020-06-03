package storage

import (
	"encoding/binary"
	"kvartalochain/common"

	"github.com/dgraph-io/badger"
)

var PREFIXNONCE = []byte("nonce")
var PREFIXHISTORY = []byte("history")

func GetBalance(db *badger.DB, addr common.Address) (uint64, error) {
	var balance uint64
	// TODO when addr not found in db, return balance 0
	err := db.View(func(txn *badger.Txn) error {
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

func GetNonce(db *badger.DB, addr common.Address) (uint64, error) {
	var nonce uint64
	err := db.View(func(txn *badger.Txn) error {
		nonceKey := append(PREFIXNONCE, addr[:]...)
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
func GetTxCount(db *badger.DB, addr common.Address) (uint64, error) {
	var count uint64
	err := db.View(func(txn *badger.Txn) error {
		countKey := append(PREFIXHISTORY, addr[:]...)
		item, err := txn.Get(countKey)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return item.Value(func(val []byte) error {
				count = binary.LittleEndian.Uint64(val)
				return err
			})
		}
		return nil
	})
	return count, err
}

func GetTx(db *badger.DB, addr common.Address, n uint64) (*common.Tx, error) {
	var nBytes [8]byte
	binary.LittleEndian.PutUint64(nBytes[:], n)

	var txBytes []byte
	err := db.View(func(txn *badger.Txn) error {
		key := append(PREFIXHISTORY, addr[:]...)
		key = append(key, nBytes[:]...)
		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return item.Value(func(val []byte) error {
				txBytes = make([]byte, len(val))
				copy(txBytes, val)
				return err
			})
		}
		return nil
	})
	tx, err := common.TxFromBytes(txBytes)
	return tx, err
}

func GetAddressHistory(db *badger.DB, addr common.Address, n uint64) ([]common.Tx, error) {
	var txs []common.Tx
	for i := 0; i < int(n); i++ {
		tx, err := GetTx(db, addr, uint64(i))
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}
	return txs, nil
}
