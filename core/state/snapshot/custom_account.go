package snapshot

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// CustomAccount is a modified version of a state.Account, where the root is replaced
// with a byte slice. This format can be used to represent full-consensus format
// or slim-snapshot format which replaces the empty root and code hash as nil
// byte slice.
//
// CustomAccount include AccountNumber and PubKey than original Account
type CustomAccount struct {
	Nonce         uint64
	Balance       *big.Int
	Root          []byte
	CodeHash      []byte
	AccountNumber uint64
	PubKeyRLP     []byte
}

// SlimAccountRLPCustom converts a state.Account content into a slim snapshot
// version RLP encoded.
//
// SlimAccountRLPCustom add accountNumber and pubKey to encode
func SlimAccountRLPCustom(nonce uint64, balance *big.Int, root common.Hash, codehash []byte, accountNumber uint64, pubKeyRLP []byte) []byte {
	data, err := rlp.EncodeToBytes(SlimAccountCustom(nonce, balance, root, codehash, accountNumber, pubKeyRLP))
	if err != nil {
		panic(err)
	}
	return data
}

// SlimAccountCustom converts a state.Account content into a slim snapshot account
//
// SlimAccountCustom add accountNumber and pubKey
func SlimAccountCustom(nonce uint64, balance *big.Int, root common.Hash, codehash []byte, accountNumber uint64, pubKeyRLP []byte) CustomAccount {
	slim := CustomAccount{
		Nonce:         nonce,
		Balance:       balance,
		AccountNumber: accountNumber,
		PubKeyRLP:     pubKeyRLP,
	}
	if root != emptyRoot {
		slim.Root = root[:]
	}
	if !bytes.Equal(codehash, emptyCode[:]) {
		slim.CodeHash = codehash
	}
	return slim
}
