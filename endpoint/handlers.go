package endpoint

import (
	"fmt"
	"kvartalochain/common"

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

func handleGetNonce(c *gin.Context) {
	addrStr := c.Param("addr")

	addr, err := common.AddressFromString(addrStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	fmt.Println("get nonce addr", addr, addr.String())
	nonce, err := getNonce(addr)
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
