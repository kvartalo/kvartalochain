package common

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/ed25519"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
)

type Address [32]byte

type PublicKey []byte

type PrivateKey ed25519.PrivateKey

/*
	the PublicKey & PrivateKey are basically wrappers over golang ed25519,
	using as Address the blake2b hash, with base58 to string representation
*/

func NewKey(r io.Reader) (PublicKey, PrivateKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	return PublicKey(pk), PrivateKey(sk), err
}

func ImportKey(sk []byte) PrivateKey {
	return PrivateKey(sk)
}
func ImportKeyString(skStr string) PrivateKey {
	return PrivateKey(base58.Decode(skStr))
}

func (sk PrivateKey) Public() PublicKey {
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, sk[32:])
	return PublicKey(publicKey)
}

func (sk PrivateKey) Sign(m []byte) []byte {
	return ed25519.Sign(ed25519.PrivateKey(sk), m)
}

func (sk PrivateKey) SignTx(tx *Tx) {
	txBytes := tx.Bytes()
	h := blake2b.Sum256(txBytes[:])
	sig := ed25519.Sign(ed25519.PrivateKey(sk), h[:])
	tx.Signature = sig
}

func VerifySignature(pk PublicKey, msg, sig []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pk), msg, sig)
}
func VerifySignatureTx(pk PublicKey, tx *Tx) bool {
	txToHash := tx.Clone()
	txToHash.Signature = []byte{}
	txBytes := txToHash.Bytes()
	h := blake2b.Sum256(txBytes[:])
	return ed25519.Verify(ed25519.PublicKey(pk), h[:], tx.Signature)
}

func (pk *PublicKey) String() string {
	return base58.Encode([]byte(*pk))
}

func (pk *PublicKey) Address() Address {
	return Address(blake2b.Sum256(*pk))
}

func (a Address) String() string {
	return base58.Encode(a[:])
}

func AddressFromString(s string) (Address, error) {
	var addrBytes [32]byte
	decodedS := base58.Decode(s)
	var empty [32]byte
	if bytes.Equal(decodedS, empty[:]) || len(decodedS) != 32 {
		return Address{}, fmt.Errorf("AddressFromString error")
	}
	copy(addrBytes[:], decodedS)
	return Address(addrBytes), nil
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Address) UnmarshalText(data []byte) error {
	addr, err := AddressFromString(string(data))
	if err != nil {
		return err
	}
	*a = addr
	return nil
}

type Tx struct {
	From      Address `json:"from" binding:"required"`
	To        Address `json:"to" binding:"required"`
	Amount    uint64  `json:"amount" binding:"required"`
	Nonce     uint64  `json:"nonce" binding:"required"`
	Signature []byte  `json:"signature" binding:"required"`
}

func NewTx(from, to Address, amount, nonce uint64) *Tx {
	return &Tx{
		From:      from,
		To:        to,
		Amount:    amount,
		Nonce:     nonce,
		Signature: []byte{},
	}
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	b := tx.Bytes()
	return json.Marshal(base58.Encode(b[:]))
}
func (tx *Tx) UnmarshalJSON(data []byte) error {
	var err error
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		panic(err)
	}
	d := base58.Decode(s)
	txB, err := TxFromBytes(d)

	tx.From = txB.From
	tx.To = txB.To
	tx.Amount = txB.Amount
	tx.Nonce = txB.Nonce
	tx.Signature = txB.Signature
	return err
}

func (tx *Tx) Bytes() [144]byte {
	var b [144]byte
	var amount [8]byte
	binary.LittleEndian.PutUint64(amount[:], tx.Amount)
	var nonce [8]byte
	binary.LittleEndian.PutUint64(nonce[:], tx.Nonce)
	copy(b[:32], tx.From[:32])
	copy(b[32:64], tx.To[:32])
	copy(b[64:72], amount[:8])
	copy(b[72:80], nonce[:8])
	copy(b[80:144], tx.Signature[:])
	return b
}

func TxFromBytes(b []byte) (*Tx, error) {
	if len(b) != 144 {
		return nil, errors.New("invalid length")
	}
	amount := binary.LittleEndian.Uint64(b[64:72])
	nonce := binary.LittleEndian.Uint64(b[72:80])
	var from, to [32]byte
	copy(from[:], b[:32])
	copy(to[:], b[32:64])
	return &Tx{
		From:      Address(from),
		To:        Address(to),
		Amount:    amount,
		Nonce:     nonce,
		Signature: b[80:144],
	}, nil
}

func (tx *Tx) Clone() *Tx {
	return &Tx{
		From:      tx.From,
		To:        tx.To,
		Amount:    tx.Amount,
		Nonce:     tx.Nonce,
		Signature: tx.Signature,
	}
}
