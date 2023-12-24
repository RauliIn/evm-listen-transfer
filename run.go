package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"transfer/config"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"

	"time"
)

var configFile = flag.String("f", "etc/deploy.yaml", "the config file")

func main() {
	//读取配置文件
	flag.Parse()

	var c config.Config
	var zeroValue = big.NewInt(0)
	gasLimit := 22000
	conf.MustLoad(*configFile, &c)
	logx.Infof("===============read config success===============")

	fromAddress, err := getAddress(c.PriKey)
	if err != nil {
		return
	}
	toAddress := common.HexToAddress(c.MintConf.ToAddr)
	//获取发送数据
	inputData := c.MintConf.InputData

	//获取rpc客户端
	client, err := ethclient.Dial(c.EthRpcConf.Url)
	if err != nil {
		logx.Errorf(" ethclient.Dial:%s", err.Error())
		return
	}
	//获取chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		logx.Errorf(" client.NetworkID:%s", err.Error())
	}

	privateKey, err := crypto.HexToECDSA(c.PriKey)

	//定时执行
	ticker := time.NewTicker(time.Duration(c.EthRpcConf.IntervalTime) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		//获取余额
		balance, err := client.PendingBalanceAt(context.Background(), *fromAddress)
		if err != nil {
			logx.Errorf("SendTransaction:insufficient funds :%s", err.Error())
			continue
		}
		if balance.Cmp(zeroValue) < 1 {
			logx.Errorf("invalid balance :%d", balance)
			continue
		}
		gasPrice, err := client.SuggestGasPrice(context.Background())

		if err != nil {
			logx.Errorf("client.SuggestGasPrice:%s", err.Error())
			continue
		}
		needGas := big.NewInt(gasPrice.Int64() * int64(gasLimit))

		value := balance.Sub(balance, needGas)
		if value.Cmp(zeroValue) < 1 {
			logx.Errorf("PendingBalanceAt :%s", err.Error())
			continue
		}
		//获取nonce
		nonce, err := client.PendingNonceAt(context.Background(), *fromAddress)
		if err != nil {
			logx.Errorf("crypto.PubkeyToAddress:%s", err.Error())
		}

		//组装交易数据
		tx := types.NewTransaction(nonce, toAddress, value, uint64(gasLimit), gasPrice, []byte(inputData))

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			logx.Errorf(" client.NetworkID:%s", err.Error())
			continue
		}
		//发送交易
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			logx.Errorf(" SendTransaction:%s", err.Error())
			continue
		}

		logx.Infof("fromAddress:%s,toAddress:%s,nonce:%d,gas:%v,txhash:%v", fromAddress, toAddress, nonce, gasPrice, signedTx.Hash())

		nonce++
	}

}

// 根据私钥获取地址
func getAddress(priKey string) (*common.Address, error) {
	privateKey, err := crypto.HexToECDSA(priKey)
	if err != nil {
		logx.Errorf("crypto.HexToECDSA:%s", err.Error())
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logx.Errorf("ecdsa.PublicKey:%s", err.Error())
		return nil, err
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, nil
}
