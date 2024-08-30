package main

import (
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	QueueSize = 1000000

	VALUE    = 1
	GasLimit = 21000
	GasPrice = 6_000_000_000 // 6gwei
	ChainId  = 666666
	ToAddr   = "0xAe95d8DA9244C37CaC0a3e16BA966a8e852Bb6D6"

	TxNumberPerAccount = 1
)

var (
	value     = big.NewInt(VALUE)
	gasLimit  = uint64(GasLimit)
	gasPrice  = big.NewInt(GasPrice)
	chainId   = big.NewInt(ChainId)
	toAddress = common.HexToAddress(ToAddr)
)

func startTxMaker(accounts []*Account, txQueue chan *types.Transaction) {
	var data []byte
	legacyTx := &types.LegacyTx{
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	}

	signer := types.NewEIP155Signer(chainId)

	for _, account := range accounts {
		for i := 0; i < TxNumberPerAccount; i++ {
			legacyTx.Nonce = account.nonce
			tx := types.NewTx(legacyTx)
			signedTx, err := types.SignTx(tx, signer, account.pk)
			if err != nil {
				log.Fatal(err)
			}
			txQueue <- signedTx
			account.nonce += 1
		}
	}

	close(txQueue)
}

func main() {
	accounts := Load()
	log.Println("load accounts finish")

	SyncNonce(accounts)
	log.Println("sync accounts nonce finish")

	wg := &sync.WaitGroup{}
	txQueue := make(chan *types.Transaction, QueueSize)

	StartTxSenderPool(wg, txQueue)
	startTxMaker(accounts, txQueue)

	log.Println("wait finish...")
	wg.Wait()
	log.Println("finish")
}
