package main

import (
	"errors"
	// "encoding/json"
	"crypto/ecdsa"
	"context"
	"fmt"
	// "io/ioutil"
	// "log"
	// "net/http"
	// "strconv"
	"strings"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/accounts/abi"

)

const transfuncode string = "0xa9059cbb"

var (
	ether       = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	// 裸秘钥 crypto.GenerateKey()
	middleKey, _    = crypto.HexToECDSA("9cc21da5ebde188b05dc6cc8727fa4e81eca2ed7d89baf872be658462e17c18a")
	middlePublicKey = "0x5efE69F7fdAfd1e4fD93CE9465A8b05d876a99c8"
	contractaddress = "0xfD98557e79ad2bd14ADdb6b5d12C1526C2E7A318"
	to = "0x6a7b119df3042f6a77ee3c1710ca54817ec41036"
	ctxtmout = 10		// rpc超时时间,默认为秒
	cpgDecimals = 18
	gasLimit uint64 = 378000
	priChainID int64 = 15 	// main;test;private表示链类型
	cpgABI = `[
		{
			"constant": false,
			"inputs": [
				{
					"name": "_spender",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "approve",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "totalSupply",
			"outputs": [
				{
					"name": "",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_from",
					"type": "address"
				},
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transferFrom",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "burn",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
				{
					"name": "_owner",
					"type": "address"
				}
			],
			"name": "balanceOf",
			"outputs": [
				{
					"name": "balance",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transfer",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
				{
					"name": "_owner",
					"type": "address"
				},
				{
					"name": "_spender",
					"type": "address"
				}
			],
			"name": "allowance",
			"outputs": [
				{
					"name": "remaining",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "burner",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Burn",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "from",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "to",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Transfer",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "owner",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "spender",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Approval",
			"type": "event"
		}
	]`
)
////////////////////////////common////////////////////////////////////////////
func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}
func to_32byte_str(para string) (str string,ret bool) {
	maxsum := 32 * 2
	strlen := len(para)
	if strlen > maxsum {
		return "",false
	}
	str = strings.Repeat("0",maxsum-strlen)
	str += para
	return str,true
}
func make_contract_abi(strs... string) (str string,ret bool){
	str = transfuncode
	ret = true
	for _, s := range strs {
		ss,ok := to_32byte_str(s)
		if !ok {
			ret = false
			return "",ret
		}
		str += ss
	}
	return str,ret
}
// make transfer code
func makeTransferCode(toAddress string, amount *big.Int) ([]byte,error){
	sval := toHexInt(amount)
	to := toAddress
	if strings.HasPrefix(toAddress,"0x") || strings.HasPrefix(toAddress,"0X") {
		to = toAddress[2:]	// 去掉目标地址0x
	}
	sdata,ok := make_contract_abi(to,sval)
	if ok {
		data := common.FromHex(sdata)
		return data,nil
	}
	return nil,errors.New("param was too long")
}
func getChainID(chainType string) *big.Int {
	chainID := big.NewInt(0)
	if strings.Compare(chainType,"main") == 0 {
		chainID = params.TestnetChainConfig.ChainId
	}else if strings.Compare(chainType,"test") == 0 {
		chainID = params.TestnetChainConfig.ChainId
	}else if strings.Compare(chainType,"private") == 0 {
		chainID = big.NewInt(priChainID)
	}else {
		fmt.Println("invalid chain type")
	}
	return chainID
}
////////////////////////////common////////////////////////////////////////////

////////////////////////////interface////////////////////////////////////////////
func getBalance(client *ethclient.Client, addr string) *big.Int {
	balance, _ := client.BalanceAt(context.Background(), common.HexToAddress(addr), nil)
	balance.Div(balance, ether)
	return balance
}
func sendETHCoin(client *ethclient.Client, key *ecdsa.PrivateKey, to string,
	amount *big.Int,chainType string) (string,error) {
	chainID := getChainID(chainType)
	nonce, _ := client.PendingNonceAt(context.Background(), common.HexToAddress(middlePublicKey))
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	tx := types.NewTransaction(nonce, common.HexToAddress(to), amount, 37800, gasPrice, nil)
	signed, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), key)

	ret := client.SendTransaction(context.Background(), signed)
	txhash := tx.Hash().String()
	return txhash,ret
}
func sendTokenCoin(client *ethclient.Client, key *ecdsa.PrivateKey,contractaddress string,
	to string, amount *big.Int,chainType string) (string,error) {

	d := time.Now().Add(time.Duration(ctxtmout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	data,dok := makeTransferCode(to,amount)
	if dok != nil {
		return "",dok
	}
	chainID := getChainID(chainType)
	nonce, _ := client.PendingNonceAt(context.Background(), common.HexToAddress(middlePublicKey))
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	tx := types.NewTransaction(nonce, common.HexToAddress(contractaddress), nil, gasLimit, gasPrice, data)
	signed, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	ret := client.SendTransaction(ctx, signed)

	txhash := tx.Hash().String()
	return txhash,ret
}
func sendTokenCoinForABI(client *ethclient.Client, key *ecdsa.PrivateKey,contractaddress string,
	to string, amount *big.Int,chainType string) (string,error) {

	d := time.Now().Add(time.Duration(ctxtmout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	cpgabi, err := abi.JSON(strings.NewReader(cpgABI))
	if err != nil {
		fmt.Println("make abi failed,err=",err)
		return "",err
	}
	bytesData, berr := cpgabi.Pack("transfer", common.HexToAddress(to), amount)
	if berr != nil {
		fmt.Println("abi pack failed,err=",berr)
		return "",berr
	}
	chainID := getChainID(chainType)
	nonce, _ := client.PendingNonceAt(context.Background(), common.HexToAddress(middlePublicKey))
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	tx := types.NewTransaction(nonce, common.HexToAddress(contractaddress), nil, gasLimit, gasPrice, bytesData)
	signed, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	ret := client.SendTransaction(ctx, signed)
	txhash := tx.Hash().String()
	return txhash,ret
}
////////////////////////////interface/////////////////////////////////////////////


////////////////////////////test//////////////////////////////////////////////////
func testSendETHCoin(client *ethclient.Client){
	amount := big.NewInt(40)
	amount.Mul(amount, ether)
	chainType := "private"
	txhash,err := sendETHCoin(client,middleKey,to,amount,chainType)
	fmt.Println("err=",err,";txhash=",txhash)
}
func testSendTokenCoin1(client *ethclient.Client){
	amount := big.NewInt(2000)
	dd := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(cpgDecimals)), nil)
	amount.Mul(amount, dd)
	chainType := "private"
	txhash,err:=sendTokenCoin(client,middleKey,contractaddress,to,amount,chainType)
	fmt.Println("err=",err,";txhash=",txhash)
}
func testSendTokenCoin2(client *ethclient.Client){
	amount := big.NewInt(2000)
	dd := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(cpgDecimals)), nil)
	amount.Mul(amount, dd)
	chainType := "private"
	txhash,err:=sendTokenCoinForABI(client,middleKey,contractaddress,to,amount,chainType)
	fmt.Println("err=",err,";txhash=",txhash)
}
////////////////////////////test//////////////////////////////////////////////////

func main() {

	client, _ := ethclient.Dial("http://127.0.0.1:8545")

	testSendTokenCoin1(client)

	// b := getBalance(client, middlePublicKey)
	// fmt.Println(b)
}
