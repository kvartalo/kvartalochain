package endpoint

import (
	"github.com/dgraph-io/badger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db *badger.DB

func newApiService() *gin.Engine {
	api := gin.Default()
	api.Use(cors.Default())
	api.GET("/info", handleInfo)
	api.GET("/balance/:addr", handleGetBalance)
	api.GET("/nonce/:addr", handleGetNonce)
	return api
}

func Serve(badgerdb *badger.DB) *gin.Engine {
	db = badgerdb
	return newApiService()
}
