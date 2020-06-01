package chain

import (
	"fmt"

	"github.com/dgraph-io/badger"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type KvartaloABCI struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

var _ abcitypes.Application = (*KvartaloABCI)(nil)

func NewKvartaloApplication(db *badger.DB) *KvartaloABCI {
	return &KvartaloABCI{
		db: db,
	}
}

func (KvartaloABCI) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (KvartaloABCI) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *KvartaloABCI) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		fmt.Println("CheckTx not valid, code: ", code)
	}
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

func (app *KvartaloABCI) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		fmt.Println("CheckTx not valid, code: ", code)
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	code = app.performTx(req.Tx)
	if code != 0 {
		fmt.Println("code", code)
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *KvartaloABCI) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *KvartaloABCI) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data

	// TODO app.getBalance

	return
}

func (app *KvartaloABCI) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *KvartaloABCI) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (KvartaloABCI) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}
