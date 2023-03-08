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
// CustomAccount include Extra and original Account
type CustomAccount struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte
	CodeHash []byte
	Extra    []byte
}

// SlimAccountRLPCustom converts a state.Account content into a slim snapshot
// version RLP encoded.
//
// SlimAccountRLPCustom add accountNumber and pubKey to encode
func SlimAccountRLPCustom(nonce uint64, balance *big.Int, root common.Hash, codehash []byte, extra []byte) []byte {
	data, err := rlp.EncodeToBytes(SlimAccountCustom(nonce, balance, root, codehash, extra))
	if err != nil {
		panic(err)
	}
	return data
}

// SlimAccountCustom converts a state.Account content into a slim snapshot account
//
// SlimAccountCustom add extra
func SlimAccountCustom(nonce uint64, balance *big.Int, root common.Hash, codehash []byte, extra []byte) CustomAccount {
	slim := CustomAccount{
		Nonce:   nonce,
		Balance: balance,
		Extra:   extra,
	}
	if root != emptyRoot {
		slim.Root = root[:]
	}
	if !bytes.Equal(codehash, emptyCode[:]) {
		slim.CodeHash = codehash
	}
	return slim
}
