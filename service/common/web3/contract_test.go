package web3

import (
	"fmt"
	"math/big"
	"testing"
	"wtc/common/types"
	"wtc/conf"
)

func TestContract(t *testing.T) {
	contract, err := NewContract(conf.ChainUrl, types.Address(conf.HexAddr), conf.HexKey)
	if err != nil {
		t.Errorf("NewContract() error = %v", err)
		return
	}
	t.Log(fmt.Sprintf("%064x", 32))
	tx, err := contract.MintNFT("0x1111111111111111111111111111111111111111", "abcffffffffgffdgdfgdfgdfsgdgfdgs")
	if err != nil {
		t.Errorf("Contract.MintNFT() error = %v", err)
		return
	}
	t.Log(tx)
	tx, err = contract.TransferNFT("0x1111111111111111111111111111111111111111", "0x1111111111111111111111111111111111111110", big.NewInt(9))
	if err != nil {
		t.Errorf("Contract.TransferNFT() error = %v", err)
		return
	}
	t.Log(tx)
	events, err := contract.LatestEvents()
	if err != nil {
		t.Errorf("Contract.LatestEvents() error = %v", err)
		return
	}
	t.Log(events)
}
