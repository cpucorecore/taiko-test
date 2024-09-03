package main

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	QueueSize = 200000
)

const (
	VALUE    = 1
	GasLimit = 21000
	ChainId  = 666666
	ToAddr   = "0xAe95d8DA9244C37CaC0a3e16BA966a8e852Bb6D6"
)

// GasPrice for Legacy tx
const (
	GasPrice = 6_000_000_000 // 6gwei
)

// MaxPriorityFeePerGas for 1559 tx
const (
	MaxPriorityFeePerGas = 6_000_000_000
	BaseFeePerGas        = 1 // should query from the network
)

const (
	TxNumberPerAccount = 5
)

var (
	chainId   = big.NewInt(ChainId)
	toAddress = common.HexToAddress(ToAddr)
	value     = big.NewInt(VALUE)
	gasLimit  = uint64(GasLimit)
)

var (
	gasPrice = big.NewInt(GasPrice)
)

var (
	maxPriorityFeePerGas = big.NewInt(MaxPriorityFeePerGas)
	maxFeePerGas         = big.NewInt(BaseFeePerGas + MaxPriorityFeePerGas)
)

func startLegacyTxMaker(accounts []*Account, txQueue chan *types.Transaction) {
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

func start1559TxMaker(accounts []*Account, txQueue chan *AccountTransaction) {
	var data []byte

	dynamicFeeTx := &types.DynamicFeeTx{
		Gas:       gasLimit,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	}

	signer := types.NewLondonSigner(chainId)

	for _, account := range accounts {
		for i := 0; i < TxNumberPerAccount; i++ {
			dynamicFeeTx.Nonce = account.nonce
			tx := types.NewTx(dynamicFeeTx)
			signedTx, err := types.SignTx(tx, signer, account.pk)
			if err != nil {
				log.Fatal(err)
			}
			txQueue <- &AccountTransaction{
				tx:   signedTx,
				addr: account.addr,
			}
			account.nonce += 1
		}
	}

	close(txQueue)
}
