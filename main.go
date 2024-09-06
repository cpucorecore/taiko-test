package main

import (
	"log"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type AccountTransaction struct {
	tx   *types.Transaction
	addr common.Address
}

var bar *pb.ProgressBar

func main() {
	accounts := Load()
	log.Println("load accounts finish")

	SyncNonce(accounts)
	log.Println("sync accounts nonce finish")

	wg := &sync.WaitGroup{}
	txQueue := make(chan *AccountTransaction, QueueSize)

	bar = pb.StartNew(TxNumberPerAccount * len(accounts))

	StartTxSenderPool(wg, txQueue)
	start1559TxMaker(accounts, txQueue)

	log.Println("wait finish...")
	wg.Wait()
	bar.Finish()
	log.Println("finish")
}
