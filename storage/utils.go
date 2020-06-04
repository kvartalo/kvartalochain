package storage

import (
	"encoding/binary"
	"kvartalochain/common"

	"github.com/dgraph-io/badger"
)

var PREFIXNONCE = []byte("nonce")
var PREFIXHISTORY = []byte("history")

func GetBalance(db *Storage, addr common.Address) uint64 {
	balanceBytes := db.Get(addr[:])
	if len(balanceBytes) == 0 {
		return uint64(0)
	}
	return binary.LittleEndian.Uint64(balanceBytes)
}

func GetNonce(db *Storage, addr common.Address) uint64 {
	nonceKey := append(PREFIXNONCE, addr[:]...)
	nonceBytes := db.Get(nonceKey)
	if len(nonceBytes) == 0 {
		return uint64(0)
	}
	return binary.LittleEndian.Uint64(nonceBytes)
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
