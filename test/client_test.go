package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"kvartalochain/common"
	"kvartalochain/endpoint"
	"kvartalochain/storage"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

var nodeurl = "http://127.0.0.1:26657"

// var nodeurl = "http://127.0.0.1:5000"

var db *badger.DB

var clienttest bool

var initBalance *bool
var oppositeDirection *bool
var addBalance *bool
var addrFlag *string
var amountFlag *int

func init() {
	initBalance = flag.Bool("initBalance", false, "init balance")
	oppositeDirection = flag.Bool("oppositeDirection", false, "tx in oposite direction")
	addBalance = flag.Bool("addBalance", false, "Add balance to address")
	addrFlag = flag.String("addr", "DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN", "Address to add balance")
	amountFlag = flag.Int("amount", 0, "Amount to be added")
	if os.Getenv("CLIENT") == "test" {
		clienttest = true
	}
	if !clienttest {
		return
	}
}

func setDbBalance(db *storage.Storage, addr common.Address, balance uint64) {
	fmt.Println("ADD BALANCE")
	var balanceBytes [8]byte
	binary.LittleEndian.PutUint64(balanceBytes[:], balance)
	db.Set(addr[:], balanceBytes[:])
	// txn := db.NewTransaction(true)
	// if err := txn.Set(addr[:], balanceBytes[:]); err == badger.ErrTxnTooBig {
	//         _ = txn.Commit()
	// }
	// _ = txn.Commit()
}

func printDbBalance(db *storage.Storage, addr common.Address) {
	// view if balance is updated
	balance := storage.GetBalance(db, addr)
	// var balance uint64
	// err := db.View(func(txn *badger.Txn) error {
	//         item, err := txn.Get(addr[:])
	//         if err != nil && err != badger.ErrKeyNotFound {
	//                 return err
	//         }
	//         if err == nil {
	//                 return item.Value(func(val []byte) error {
	//                         balance = binary.LittleEndian.Uint64(val)
	//                         return err
	//                 })
	//         }
	//         return nil
	// })
	// if err != nil {
	//         panic(err)
	// }
	fmt.Println("Address:", addr, ", Balance:", balance)
}

func getBalance(addr common.Address) (uint64, error) {
	apiurl := "http://127.0.0.1:3000"
	r, err := http.Get(apiurl + "/balance/" + addr.String())
	if err != nil {
		return 0, err
	}
	var msg endpoint.GetBalanceMsg
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		return 0, err
	}
	fmt.Println("getBalance", msg)

	return msg.Balance, nil
}

func TestClient(t *testing.T) {
	if !clienttest {
		t.Skip()
	}
	flag.Parse()

	if *addBalance {
		db, err := storage.NewStorage("../data")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open storage db: %v", err)
			os.Exit(1)
		}

		archiveDb, err := badger.Open(badger.DefaultOptions("../data"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
			os.Exit(1)
		}
		defer archiveDb.Close()

		fmt.Printf("Add Balance %v to addr %s\n", *amountFlag, *addrFlag)
		addr, err := common.AddressFromString(*addrFlag)
		if err != nil {
			panic(err)
		}
		setDbBalance(db, addr, uint64(*amountFlag))
		db.Commit()
		printDbBalance(db, addr)

		os.Exit(0)
	}

	var skStr0, skStr1 string
	a := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
	b := "8h3u7NfgvUJsHJgKDUKwwVL1iZd3cwRtntpTfJ5Mefz2"
	if !*oppositeDirection {
		skStr0 = a
		skStr1 = b
	} else {
		skStr0 = b
		skStr1 = a
	}

	sk0 := common.ImportKeyString(skStr0)
	pk0 := sk0.Public()
	addr0 := pk0.Address()
	assert.Equal(t, "DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN", addr0.String())

	addr0FromStr, err := common.AddressFromString("DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN")
	if err != nil {
		panic(err)
	}
	fmt.Println(addr0FromStr.String())
	fmt.Println(addr0.String())
	assert.Equal(t, addr0FromStr, addr0)

	sk1 := common.ImportKeyString(skStr1)
	pk1 := sk1.Public()
	addr1 := pk1.Address()
	assert.Equal(t, "HzeXxgjb589tVBs991jAyLUX7wreSZvrWnRxdGQS4co2", addr1.String())

	if *initBalance {
		db, err := storage.NewStorage("../data")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open storage db: %v", err)
			os.Exit(1)
		}
		archiveDb, err := badger.Open(badger.DefaultOptions("../data"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
			os.Exit(1)
		}
		defer archiveDb.Close()
		fmt.Println("initBalance flag activated")
		setDbBalance(db, addr0, 100)
		setDbBalance(db, addr1, 0)
		db.Commit()
		printDbBalance(db, addr0)
		printDbBalance(db, addr1)
		os.Exit(0)
	}

	var bal uint64
	bal, err = getBalance(addr0)
	assert.Nil(t, err)
	fmt.Println("b0", bal)
	assert.Equal(t, uint64(100), bal)
	bal, err = getBalance(addr1)
	assert.Nil(t, err)
	fmt.Println("b1", bal)
	assert.Equal(t, uint64(0), bal)

	// create and sign tx
	tx := &common.Tx{
		From:      addr0,
		To:        addr1,
		Amount:    10,
		Signature: []byte{},
	}
	sk0.SignTx(tx)
	assert.NotEqual(t, []byte{}, tx.Signature)

	assert.True(t, common.VerifySignatureTx(tx))

	// send tx
	txHex := hex.EncodeToString(tx.Bytes())

	fmt.Println("sending", nodeurl+`/broadcast_tx_commit?tx="`+txHex+`"`)
	resp, err := http.Get(nodeurl + `/broadcast_tx_commit?tx="` + txHex + `"`)
	assert.Nil(t, err)
	fmt.Println(resp)

	fmt.Println("waiting some blocks")
	time.Sleep(2 * time.Second)
	// print balances after tx
	bal, err = getBalance(addr0)
	assert.Nil(t, err)
	fmt.Println("b0", bal)
	assert.Equal(t, uint64(90), bal)
	bal, err = getBalance(addr1)
	assert.Nil(t, err)
	fmt.Println("b1", bal)
	assert.Equal(t, uint64(10), bal)
}
