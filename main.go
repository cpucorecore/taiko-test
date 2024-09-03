package main

import (
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type AccountTransaction struct {
	tx   *types.Transaction
	addr common.Address
}

func main() {
	accounts := Load()
	log.Println("load accounts finish")

	SyncNonce(accounts)
	log.Println("sync accounts nonce finish")

	wg := &sync.WaitGroup{}
	txQueue := make(chan *AccountTransaction, QueueSize)

	StartTxSenderPool(wg, txQueue)
	start1559TxMaker(accounts, txQueue)

	log.Println("wait finish...")
	wg.Wait()
	log.Println("finish")
}
