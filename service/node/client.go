package node

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/rpc"
	"wtc/common/types"
	"wtc/common/utils"
)

var NotFound = fmt.Errorf("not found")

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	*rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	rpc, err := rpc.Dial(rawurl)
	return &Client{rpc}, err
}

type Big big.Int

func (b *Big) UnmarshalJSON(input []byte) error {
	return (*big.Int)(b).UnmarshalJSON(input[1 : len(input)-1])
}

func (c *Client) ChainId(ctx context.Context) (*big.Int, error) {
	var hex Big
	if err := c.CallContext(ctx, &hex, "eth_chainId"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *Client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var hex Big
	if err := c.CallContext(ctx, &hex, "eth_blockNumber"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex Big
	if err := c.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *Client) EstimateGas(ctx context.Context, msg map[string]interface{}) (uint64, error) {
	var result string
	err := c.CallContext(ctx, &result, "eth_estimateGas", msg)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(result[2:], 16, 64)
}

func (c *Client) PendingNonceAt(ctx context.Context, account types.Address) (uint64, error) {
	var result string
	err := c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(result[2:], 16, 64)
}

func (c *Client) TxParams(ctx context.Context, addr types.Address) (chainId, gasPrice *big.Int, nonce uint64, err error) {
	chainId, err = c.ChainId(ctx)
	if err != nil {
		return
	}
	gasPrice, err = c.SuggestGasPrice(ctx)
	if err != nil {
		return
	}
	nonce, err = c.PendingNonceAt(ctx, addr)
	return
}

func (c *Client) SendTx(prv *secp256k1.PrivateKey, to *types.Address, value *big.Int, data types.Data) (result types.Hash, err error) {
	ctx := context.Background()
	addr := utils.PubkeyToAddress(prv.PubKey())
	gas, err := c.EstimateGas(ctx, map[string]interface{}{
		"from":  addr,
		"to":    to,
		"value": "0x" + value.Text(16),
		"data":  data,
	})
	if err != nil {
		return
	}
	chainId, gasPrice, nonce, err := c.TxParams(ctx, addr)
	if err != nil {
		return
	}
	tx := utils.NewTx(nonce, to, value, gas, gasPrice, data)
	rawTx, err := utils.SignTx(tx, chainId, prv)
	if err != nil {
		return
	}
	err = c.CallContext(ctx, &result, "eth_sendRawTransaction", rawTx)
	return
}

func (c *Client) Deploy(prv *secp256k1.PrivateKey, value *big.Int, data types.Data) (types.Address, error) {
	ctx := context.Background()
	addr := utils.PubkeyToAddress(prv.PubKey())
	gas, err := c.EstimateGas(ctx, map[string]interface{}{
		"from":  addr,
		"value": "0x" + value.Text(16),
		"data":  data,
	})
	if err != nil {
		return "", err
	}
	chainId, gasPrice, nonce, err := c.TxParams(ctx, addr)
	if err != nil {
		return "", err
	}
	tx := utils.NewTx(nonce, nil, value, gas, gasPrice, data)
	rawTx, err := utils.SignTx(tx, chainId, prv)
	if err != nil {
		return "", err
	}
	err = c.CallContext(ctx, nil, "eth_sendRawTransaction", rawTx)
	return utils.CreateAddress(addr, nonce), err
}

// EventLog transaction log
type EventLog struct {
	Address     types.Address `json:"address"`         //The contract address
	Topics      []types.Hash  `json:"topics"`          //topic
	Data        string        `json:"data"`            //data
	Removed     bool          `json:"removed"`         //whether to remove
	BlockNumber types.Long    `json:"blockNumber"`     //block number
	TxHash      types.Hash    `json:"transactionHash"` //The transaction hash
	Index       types.Long    `json:"logIndex"`        //The serial number in the transaction
}

func (c *Client) GetLogs(filter map[string]any) (logs []*EventLog, err error) {
	err = c.CallContext(context.Background(), &logs, "eth_getLogs", filter)
	return
}
