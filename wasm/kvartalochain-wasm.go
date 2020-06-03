package main

import (
	"encoding/hex"
	"strconv"
	"syscall/js"

	"kvartalochain/common"
)

func main() {
	c := make(chan struct{}, 0)
	println("WASM kvartalochain initialized")
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Set("newKey", js.FuncOf(newKey))
	js.Global().Set("newTxAndSign", js.FuncOf(newTxAndSign))
}

func newKey(this js.Value, values []js.Value) interface{} {
	pk, sk, err := common.NewKey()
	if err != nil {
		return js.ValueOf(err.Error())
	}

	r := make(map[string]interface{})
	r["sk"] = sk.String()
	r["pk"] = pk.String()
	r["address"] = pk.Address().String()
	return r
}

func newTxAndSign(this js.Value, values []js.Value) interface{} {
	skStr := values[0].String()
	toStr := values[1].String()
	amountStr := values[2].String()
	nonceStr := values[3].String()

	sk := common.ImportKeyString(skStr)
	to, err := common.AddressFromString(toStr)
	if err != nil {
		return js.ValueOf(err.Error())
	}
	amountInt, err := strconv.Atoi(amountStr)
	if err != nil {
		return js.ValueOf(err.Error())
	}
	amount := uint64(amountInt)
	nonceInt, err := strconv.Atoi(nonceStr)
	if err != nil {
		return js.ValueOf(err.Error())
	}
	nonce := uint64(nonceInt)

	from := sk.Public().Address()

	tx := common.NewTx(from, to, amount, nonce)
	sk.SignTx(tx)

	r := make(map[string]interface{})
	r["from"] = tx.From.String()
	r["to"] = tx.To.String()
	r["amount"] = strconv.Itoa(int(tx.Amount))
	r["nonce"] = strconv.Itoa(int(tx.Nonce))
	r["txHex"] = hex.EncodeToString(tx.Bytes())
	return r
}
