package main

import "testing"

func TestLoad(t *testing.T) {
	accounts := Load()
	for _, account := range accounts {
		t.Log(account.addr.Hex())
	}

	SyncNonce(accounts)
	t.Log(accounts[len(accounts)-1].addr.Hex())
	t.Log(accounts[len(accounts)-1].nonce)
}
