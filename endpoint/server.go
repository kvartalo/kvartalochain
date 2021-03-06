package endpoint

import (
	"kvartalochain/storage"

	"github.com/dgraph-io/badger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db *storage.Storage
var archiveDb *badger.DB

func newApiService() *gin.Engine {
	api := gin.Default()
	api.Use(cors.Default())
	api.GET("/info", handleInfo)
	api.GET("/balance/:addr", handleGetBalance)
	api.GET("/nonce/:addr", handleGetNonce)
	api.POST("/tx", handlePostTx)
	api.GET("/history/:addr", handleGetHistory)
	return api
}

func Serve(sto *storage.Storage, badgerdb *badger.DB) *gin.Engine {
	db = sto
	archiveDb = badgerdb
	return newApiService()
}
