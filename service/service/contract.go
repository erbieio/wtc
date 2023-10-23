package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"wtc/common/utils"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"wtc/common/types"
	"wtc/node"
)

// Contract 链上合约操作对象
type Contract struct {
	client    *node.Client          //Ethereum RPC client
	prv       *secp256k1.PrivateKey //Private key object
	addr      types.Address
	curHeight *big.Int
}

var batchCode = "0x608060405234801561"

func NewContract(url string, hexAddr types.Address, hexKey string) (c *Contract, err error) {
	client, err := node.Dial(url)
	if err != nil {
		return
	}
	prv, err := utils.HexToECDSA(hexKey)
	if err != nil {
		return
	}
	if hexAddr == "" {
		hexAddr, err = client.Deploy(prv, new(big.Int), types.Data(batchCode))
		if err != nil {
			return nil, err
		}
	}
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		return
	}
	c = &Contract{client, prv, hexAddr, number}
	return
}

// MintNFT 铸造NFT，amount为铸造nft的个数
func (c *Contract) MintNFT(to types.Address, uri types.Address) (types.Hash, error) {
	data := "0xe4e2a53a000000000000000000000000" + types.Data(to[2:])
	data += "0000000000000000000000000000000000000000000000000000000000000040"
	data += types.Data(fmt.Sprintf("%064x", len(uri)))
	data += types.Data(hex.EncodeToString([]byte(uri)))
	if padding := len(uri) % 32; padding > 0 {
		data += types.Data("0000000000000000000000000000000000000000000000000000000000000000")[2*padding:]
	}
	return c.client.SendTx(c.prv, &c.addr, new(big.Int), data)
}

// TransferNFT 转移NFT
func (c *Contract) TransferNFT(from, to types.Address, id *big.Int) (types.Hash, error) {
	data := "0x905bd5e4000000000000000000000000" + types.Data(from[2:])
	data += "000000000000000000000000" + types.Data(to[2:])
	data += types.Data(fmt.Sprintf("%064x", id.Uint64())) //todo: not support more than 2^64-1
	return c.client.SendTx(c.prv, &c.addr, new(big.Int), data)
}

type Event struct {
	BlockNumber int64
	TxHash      types.Hash
	From        types.Address //if mint, this value is 0
	To          types.Address
	Id          *big.Int
}

// LatestEvents 最新的交易结果（上一次调用的高度到最新区块之间的）
func (c *Contract) LatestEvents() (events []*Event, err error) {
	newHeight, err := c.client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	// topics: Keccak256Hash([]byte("Transfer(address,address,uint256)")
	filter := map[string]any{
		"address":   c.addr,
		"topics":    []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},
		"fromBlock": "0x" + c.curHeight.Text(16),
		"toBlock":   "0x" + newHeight.Text(16),
	}
	logs, err := c.client.GetLogs(filter)
	if err != nil {
		return
	}
	c.curHeight.SetInt64(newHeight.Int64() + 1)
	for _, log := range logs {
		events = append(events, &Event{
			BlockNumber: int64(log.BlockNumber),
			TxHash:      log.TxHash,
			From:        types.Address("0x" + log.Topics[1][26:]),
			To:          types.Address("0x" + log.Topics[2][26:]),
			Id:          utils.HexToBigInt(string(log.Topics[3][2:])),
		})
	}
	return
}
