# kvartalochain [![GoDoc](https://godoc.org/github.com/kvartalo/kvartalochain?status.svg)](https://godoc.org/github.com/kvartalo/kvartalochain) [![Go Report Card](https://goreportcard.com/badge/github.com/kvartalo/kvartalochain)](https://goreportcard.com/report/github.com/kvartalo/kvartalochain) [![Test](https://github.com/kvartalo/kvartalochain/workflows/Test/badge.svg)](https://github.com/kvartalo/kvartalochain/actions?query=workflow%3ATest)

[kvartalo](https://kvartalo.xyz) chain node, using [Tendermint](https://tendermint.com).

This repo is a work in progress.


## Details
- Blockchain based on `Tendermint Core`
- Keys & Signatures using `btcec` https://godoc.org/github.com/btcsuite/btcd/btcec
- Address -> Hash `blake2b` of the `PublicKey` with a `nonce`, encoded in `base58`

## Test
- unit test:
```
go test ./...
```

- int test
```
# run the node in one terminal
go run main.go --config ~/path/to/tendermint/config/config.toml start

# run the client test
cd test
CLIENT=test go test
```
