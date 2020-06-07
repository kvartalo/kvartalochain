package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
)

type PrivateKey struct{ *btcec.PrivateKey }
type PublicKey struct{ *btcec.PublicKey }
type Address [32]byte

/*
	the PublicKey & PrivateKey are basically wrappers over golang btcec
	(https://godoc.org/github.com/btcsuite/btcd/btcec), using as Address
	the blake2b hash, with base58 to string representation
*/

func NewKey() (*PublicKey, *PrivateKey, error) {
	sk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, nil, err
	}
	pk := sk.PubKey()
	return &PublicKey{pk}, &PrivateKey{sk}, err
}

func ImportKey(b []byte) *PrivateKey {
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), b)
	return &PrivateKey{sk}
}
func ImportKeyString(skStr string) *PrivateKey {
	return ImportKey(base58.Decode(skStr))
}

func (sk PrivateKey) Bytes() []byte {
	return sk.Serialize()
}
func (sk PrivateKey) String() string {
	return base58.Encode(sk.Bytes())
}
func (sk PrivateKey) Public() *PublicKey {
	return &PublicKey{sk.PubKey()}
}

func (sk PrivateKey) HashAndSign(m []byte) ([]byte, error) {
	h := blake2b.Sum256(m)
	sig, err := btcec.SignCompact(btcec.S256(), sk.PrivateKey, h[:], false)
	return sig, err
}

func (sk PrivateKey) SignTx(tx *Tx) error {
	txBytes := tx.Bytes()
	sig, err := sk.HashAndSign(txBytes[:])
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

func VerifySignature(addr *Address, msg, sigB []byte) bool {
	h := blake2b.Sum256(msg)
	pkRec, _, err := btcec.RecoverCompact(btcec.S256(), sigB, h[:])
	if err != nil {
		return false
	}
	publicKeyRec := PublicKey{pkRec}
	addrRec := publicKeyRec.Address()
	if !bytes.Equal(addrRec[:], addr[:]) {
		return false
	}
	return true
}
func VerifySignatureTx(tx *Tx) bool {
	txToHash := tx.Clone()
	txToHash.Signature = []byte{}
	txBytes := txToHash.Bytes()
	h := blake2b.Sum256(txBytes)

	pkRec, _, err := btcec.RecoverCompact(btcec.S256(), tx.Signature, h[:])
	if err != nil {
		return false
	}
	publicKeyRec := PublicKey{pkRec}
	addrRec := publicKeyRec.Address()
	if !bytes.Equal(addrRec[:], tx.From[:]) {
		return false
	}
	return true
}

func (pk *PublicKey) Bytes() []byte {
	return pk.SerializeCompressed()
}

func (pk *PublicKey) String() string {
	return base58.Encode(pk.Bytes())
}

func (pk *PublicKey) Address() Address {
	return Address(blake2b.Sum256(pk.Bytes()))
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

type TxType byte

const TxTypeNormal = TxType(0)
const TxTypeMint = TxType(1)

func TxTypeFromByte(b byte) TxType {
	return TxType(b)
}

func (t TxType) String() string {
	switch t {
	case TxTypeNormal:
		return "TxTypeNormal"
	case TxTypeMint:
		return "TxTypeMint"
	default:
		return "TxTypeUndefined"

	}
}

type Tx struct {
	Type      TxType  `json:"type" binding:"required"`
	From      Address `json:"from" binding:"required"`
	To        Address `json:"to" binding:"required"`
	Amount    uint64  `json:"amount" binding:"required"`
	Nonce     uint64  `json:"nonce" binding:"required"`
	Signature []byte  `json:"signature" binding:"required"`
	// TODO timestamp (outside signature)
}

// NewTx returns a Tx data structure. By default, is a TxTypeNormal tx type.
func NewTx(from, to Address, amount, nonce uint64) *Tx {
	return &Tx{
		Type:      TxTypeNormal,
		From:      from,
		To:        to,
		Amount:    amount,
		Nonce:     nonce,
		Signature: []byte{},
	}
}

// func (tx *Tx) MarshalJSON() ([]byte, error) {
//         b := tx.Bytes()
//         return json.Marshal(base58.Encode(b[:]))
// }
// func (tx *Tx) UnmarshalJSON(data []byte) error {
//         var err error
//         var s string
//         err = json.Unmarshal(data, &s)
//         if err != nil {
//                 panic(err)
//         }
//         d := base58.Decode(s)
//         txB, err := TxFromBytes(d)
//
//         tx.From = txB.From
//         tx.To = txB.To
//         tx.Amount = txB.Amount
//         tx.Nonce = txB.Nonce
//         tx.Signature = txB.Signature
//         return err
// }

func (tx *Tx) Bytes() []byte {
	var b []byte
	var amount [8]byte
	binary.LittleEndian.PutUint64(amount[:], tx.Amount)
	var nonce [8]byte
	binary.LittleEndian.PutUint64(nonce[:], tx.Nonce)
	b = append(b, byte(tx.Type))
	b = append(b, tx.From[:32]...)
	b = append(b, tx.To[:32]...)
	b = append(b, amount[:8]...)
	b = append(b, nonce[:8]...)
	b = append(b, tx.Signature[:]...)
	return b
}

func (tx *Tx) Hex() string {
	return hex.EncodeToString(tx.Bytes())
}

func (tx *Tx) String() string {
	buf := bytes.NewBufferString("")
	fmt.Fprintf(buf, "Type: %v, ", tx.Type.String())
	fmt.Fprintf(buf, "From: %v, ", tx.From.String())
	fmt.Fprintf(buf, "To: %v, ", tx.To.String())
	fmt.Fprintf(buf, "Amount: %v, ", strconv.Itoa(int(tx.Amount)))
	fmt.Fprintf(buf, "Nonce: %v, ", strconv.Itoa(int(tx.Nonce)))
	fmt.Fprintf(buf, "Signature: %v", hex.EncodeToString(tx.Signature))
	return buf.String()
}

func TxFromBytes(b []byte) (*Tx, error) {
	if len(b) < 81 {
		return nil, fmt.Errorf("error on tx bytes format")
	}
	amount := binary.LittleEndian.Uint64(b[65:73])
	nonce := binary.LittleEndian.Uint64(b[73:81])
	var from, to [32]byte
	copy(from[:], b[1:33])
	copy(to[:], b[33:65])
	return &Tx{
		Type:      TxType(b[0]),
		From:      Address(from),
		To:        Address(to),
		Amount:    amount,
		Nonce:     nonce,
		Signature: b[81:],
	}, nil
}

func (tx *Tx) Clone() *Tx {
	return &Tx{
		Type:      tx.Type,
		From:      tx.From,
		To:        tx.To,
		Amount:    tx.Amount,
		Nonce:     tx.Nonce,
		Signature: tx.Signature,
	}
}
