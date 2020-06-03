#!/bin/sh

GOARCH=wasm GOOS=js go build -o kvartalochain.wasm kvartalochain-wasm.go
mv kvartalochain.wasm webtest/kvartalochain.wasm
cp webtest/kvartalochain.wasm ../../wallet/lib/
