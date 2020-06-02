package endpoint

import (
	"encoding/binary"
	"kvartalochain/chain"
	"kvartalochain/common"

	"github.com/dgraph-io/badger"
)

func getBalance(addr common.Address) (uint64, error) {
	var balance uint64
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

func getNonce(addr common.Address) (uint64, error) {
	var nonce uint64
	err := db.View(func(txn *badger.Txn) error {
		nonceKey := append(chain.NONCEKEY, addr[:]...)
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
