package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	TxSenderPoolSize = 100
	RPC              = "http://192.168.100.77:28545"
)

func StartTxSenderPool(wg *sync.WaitGroup, txQueue chan *AccountTransaction) {
	wg.Add(TxSenderPoolSize)
	for i := 0; i < TxSenderPoolSize; i++ {
		go startTxSender(wg, txQueue)
	}

	go func() {
		for {
			log.Printf("failCounter:%d\n", failCounter.Load())
			time.Sleep(time.Second * 1)
		}
	}()
}

const (
	MaxRetry      = 100
	RetryInterval = time.Second * 5
)

var (
	failCounter atomic.Uint64
)

func startTxSender(wg *sync.WaitGroup, txQueue chan *AccountTransaction) {
	defer wg.Done()

	c, err := ethclient.Dial(RPC)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	retryCnt := 0
	for tx := range txQueue {
		retryCnt = 0
		for retryCnt < MaxRetry {
			err = c.SendTransaction(context.Background(), tx.tx)
			if err != nil {
				failCounter.Add(1)
				retryCnt++
				if err.Error() == "already known" {
					break
				}

				if strings.Contains(err.Error(), "insufficient funds") {
					break
				}

				log.Printf("SendTransaction failed %d times, addr: %s, err: %v", retryCnt, tx.addr.Hex(), err)
				time.Sleep(RetryInterval)
			} else {
				bar.Increment()
				break
			}
		}

		if retryCnt == MaxRetry {
			txStr, _ := tx.tx.MarshalJSON()
			log.Printf("SendTransaction failed max retry times, addr: %s, tx: %s\n", tx.addr.Hex(), string(txStr))
			break
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
