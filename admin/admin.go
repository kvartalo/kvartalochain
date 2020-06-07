package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"kvartalochain/common"
	"net/http"
	"os"
)

var mint *bool
var addrFlag *string
var amountFlag *int

func main() {
	mint = flag.Bool("mint", false, "Mint coints to address")
	addrFlag = flag.String("addr", "DqF1B6iqaxeE3j4XvyPfLbba6QkQfQtwSUWBJmnQRMvN", "Address to add balance")
	amountFlag = flag.Int("amount", 0, "Amount to be added")
	flag.Parse()

	if *mint {
		addr, err := common.AddressFromString(*addrFlag)
		if err != nil {
			panic(err)
		}

		// tmp
		skStr0 := "2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG"
		sk0 := common.ImportKeyString(skStr0)
		pk0 := sk0.Public()
		addr0 := pk0.Address()

		tx := common.NewTx(addr0, addr, uint64(*amountFlag), uint64(0))
		tx.Type = common.TxTypeMint
		sk0.SignTx(tx)

		txHex := hex.EncodeToString(tx.Bytes())
		var nodeurl = "http://127.0.0.1:26657"
		fmt.Println("sending", nodeurl+`/broadcast_tx_commit?tx="`+txHex+`"`)
		resp, err := http.Get(nodeurl + `/broadcast_tx_commit?tx="` + txHex + `"`)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(resp)
	}
}
