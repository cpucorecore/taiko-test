package main

import (
	"context"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	TxSenderPoolSize = 100
	RPC              = "http://192.168.100.77:28545"
)

func StartTxSenderPool(wg *sync.WaitGroup, txQueue chan *types.Transaction) {
	wg.Add(TxSenderPoolSize)
	for i := 0; i < TxSenderPoolSize; i++ {
		go startTxSender(wg, txQueue)
	}
}

func startTxSender(wg *sync.WaitGroup, txQueue chan *types.Transaction) {
	defer wg.Done()

	c, err := ethclient.Dial(RPC)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for tx := range txQueue {
		err = c.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Printf("failed to send transaction: %v", err)
		} else {
			//j, _ := tx.MarshalJSON()
			//log.Println(string(j))
		}
	}

	log.Println("work done")
}

func SyncNonce(accounts []*Account) {
	accountChan := make(chan *Account)
	syncNonceWg := &sync.WaitGroup{}

	syncNonceWorkerNumber := 100
	syncNonceWg.Add(syncNonceWorkerNumber)
	for i := 0; i < syncNonceWorkerNumber; i++ {
		go func() {
			defer syncNonceWg.Done()
			c, err := ethclient.Dial(RPC)
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()

			for account := range accountChan {
				nonce, err := c.PendingNonceAt(context.Background(), account.addr)
				if err != nil {
					log.Fatal(err)
				}
				account.nonce = nonce
			}
		}()
	}

	for _, acc := range accounts {
		accountChan <- acc
	}
	close(accountChan)

	syncNonceWg.Wait()
}
