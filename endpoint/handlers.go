package endpoint

import (
	"fmt"
	"kvartalochain/common"
	"kvartalochain/storage"
	"net/http"

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
	balance := storage.GetBalance(db, addr)
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

func handleGetNonce(c *gin.Context) {
	addrStr := c.Param("addr")

	addr, err := common.AddressFromString(addrStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	fmt.Println("get nonce addr", addr, addr.String())
	nonce := storage.GetNonce(db, addr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"addr":  addr,
		"nonce": nonce,
	})
}

type PostTxMsg struct {
	TxHex string `json:"txHex"`
}

func handlePostTx(c *gin.Context) {
	var m PostTxMsg
	c.BindJSON(&m)
	_, err := http.Get("http://127.0.0.1:26657" + `/broadcast_tx_commit?tx="` + m.TxHex + `"`)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
	})

}

func handleGetHistory(c *gin.Context) {
	addrStr := c.Param("addr")
	addr, err := common.AddressFromString(addrStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	txCount, err := storage.GetTxCount(archiveDb, addr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	txs, err := storage.GetAddressHistory(archiveDb, addr, txCount)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"txs": txs,
	})
}
