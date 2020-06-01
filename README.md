# kvartalochain
[kvartalo](https://kvartalo.xyz) chain node, using [Tendermint](https://tendermint.com).

This repo is a work in progress.


## Details
- Blockchain based on `Tendermint Core`
- Keys & Signatures using `btcec` https://godoc.org/github.com/btcsuite/btcd/btcec
- Address -> Hash `blake2b` of the `PublicKey` with a `nonce`, encoded in `base58`

