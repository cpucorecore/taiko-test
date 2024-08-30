package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
)

type Account struct {
	pk     *ecdsa.PrivateKey
	pubKey *ecdsa.PublicKey
	addr   common.Address
	nonce  uint64
}
