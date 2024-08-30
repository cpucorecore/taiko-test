package main

import (
	"bufio"
	"crypto/ecdsa"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	f = "./scripts/priks"
)

func Load() (accounts []*Account) {
	fd, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		accounts = append(accounts, loadOne(scanner.Text()[2:]))
	}

	return accounts
}

func loadOne(hexKey string) *Account {
	pk, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}

	pubK, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	addr := crypto.PubkeyToAddress(*pubK)

	return &Account{
		pk:     pk,
		pubKey: pubK,
		addr:   addr,
		nonce:  0,
	}
}
