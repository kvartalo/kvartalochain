package endpoint

import (
	"encoding/binary"
	"fmt"
	"kvartalochain/common"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
)

type GetBalanceMsg struct {
	Addr    common.Address `json:"addr"`
	Balance uint64         `json:"balance"`
}

func handleInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func handleGetBalance(c *gin.Context) {
	addrStr := c.Param("addr")

	addr, err := common.AddressFromString(addrStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	fmt.Println("get balance addr", addr, addr.String())
	balance, err := getBalance(addr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, GetBalanceMsg{
		Addr:    addr,
		Balance: balance,
	})
}

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
