package main

import (
	"fmt"
	"testing"
)

func TestSyncNonce(t *testing.T) {
	accounts := Load()
	SyncNonce(accounts)
	for i, account := range accounts {
		fmt.Printf("%d %s nonce: %d\n", i, account.addr.Hex(), account.nonce)
	}
}
