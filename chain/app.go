package chain

import (
	"encoding/hex"
	"fmt"
	"kvartalochain/common"
	"kvartalochain/storage"

	"github.com/dgraph-io/badger"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type KvartaloABCI struct {
	archive      bool
	db           *storage.Storage // used for state, balances and nonces
	archiveDb    *badger.DB       // used for tx history archive
	currentBatch *badger.Txn
}

var _ abcitypes.Application = (*KvartaloABCI)(nil)

func NewKvartaloApplication(db *storage.Storage, archiveDb *badger.DB) *KvartaloABCI {
	return &KvartaloABCI{
		archive:   true,
		db:        db,
		archiveDb: archiveDb,
	}
}

func (KvartaloABCI) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (KvartaloABCI) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *KvartaloABCI) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	txBytes, err := hex.DecodeString(string(req.Tx))
	if err != nil {
		return abcitypes.ResponseCheckTx{Code: ERRFORMAT} // invalid tx format
	}
	tx, err := common.TxFromBytes(txBytes)
	if err != nil {
		return abcitypes.ResponseCheckTx{Code: ERRFORMAT} // invalid tx format
	}
	code := app.isValid(tx)
	if code != 0 {
		fmt.Println("CheckTx not valid, code: ", code)
		return abcitypes.ResponseCheckTx{Code: code}
	}
	// return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
	return abcitypes.ResponseCheckTx{Code: code}
}

func (app *KvartaloABCI) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.performTx(req.Tx)
	if code != 0 {
		// TODO if err, cancel tx, don't Commit()
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	return abcitypes.ResponseDeliverTx{Code: code}
}

func (app *KvartaloABCI) Commit() abcitypes.ResponseCommit {
	// TMP
	// h, err := app.db.Commit() // store chain state
	// if err != nil {
	//         fmt.Println("ERR", err)
	//         panic(err)
	// }
	app.currentBatch.Commit() // store archive history
	// fmt.Println(h)
	// return abcitypes.ResponseCommit{Data: h}

	// return abcitypes.ResponseCommit{Data: app.db.State()}
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *KvartaloABCI) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data

	// TODO app.getBalance

	return
}

func (app *KvartaloABCI) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	//
	// //
	// // TODO TMP
	// addrStr := "DR4cHJqhRM9xi6Dhx2EqpYAVejW2eDnBogabt7Zre9N4"
	// // addrStr := "DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN"
	// addr, err := common.AddressFromString(addrStr)
	// if err != nil {
	//         panic(err)
	// }
	// balance := uint64(10000)
	// var balanceBytes [8]byte
	// binary.LittleEndian.PutUint64(balanceBytes[:], balance)
	// app.db.Set(addr[:], balanceBytes[:])
	// fmt.Println("BAL", storage.GetBalance(app.db, addr))
	// // /TMP
	// //

	// app.db.Set([]byte("init"), []byte("init"))
	// app.db.Commit()

	return abcitypes.ResponseInitChain{}
}

func (app *KvartaloABCI) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.archiveDb.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (KvartaloABCI) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}
