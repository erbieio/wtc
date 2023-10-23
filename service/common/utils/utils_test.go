package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"testing"
)

func TestSignTx(t *testing.T) {
	tx := map[string]string{"data": "0x", "gas": "0x5208", "gasPrice": "0x3b9aca07", "nonce": "0x65", "to": "0x394586580ff4170c8a0244837202cbabe9070f66", "value": "0x0de0b6b3a7640000"}
	key, _ := HexToECDSA("7b2546a5d4e658d079c6b2755c6d7495edd01a686fddae010830e9c93b23e398")

	t.Logf("tx:%+v\n", tx)
	signedTx, _ := SignTx(tx, big.NewInt(51888), key)
	t.Log("signedTx:", signedTx)
}

func TestBigToAddress(t *testing.T) {
	a := big.NewInt(1)
	a.SetString("bbbbbbbbbbbbbbbbbbbbbbbbaaaaaaaaaaaaaaaaaaaaaaaa", 16)
	t.Log(a.SetString("0aabbb", 16))
	t.Log(fmt.Sprintf("%s%x", "0x1234", 0))
	t.Log(fmt.Sprintf("%s%x", "0x1234", 13))

	reg := regexp.MustCompile(`\b0[xX][\da-zA-Z]{40}\b`)
	str := reg.Find([]byte("{\"data\":[\"text\":\"Requesting 0xddA7619eE94aeD4d00eb6e72b28d96873a54B3fd \""))
	t.Log(string(str))
}

func TestKeccak256Hash(t *testing.T) {
	t.Log(Keccak256Hash([]byte("Transfer(address,address,uint256)")))
	t.Log(Keccak256Hash([]byte("Transfer(address,address,uint256)")))
	t.Log(Keccak256Hash([]byte("TransferSingle(address,address,address,uint256,uint256)")))
	t.Log(Keccak256Hash([]byte("TransferBatch(address,address,address,uint256[],uint256[])")))

	t.Log(Keccak256Hash([]byte("name()"))[:10])
	t.Log(Keccak256Hash([]byte("symbol()"))[:10])
	t.Log(Keccak256Hash([]byte("decimals()"))[:10])
	t.Log(Keccak256Hash([]byte("totalSupply()"))[:10])

	t.Log(Keccak256Hash([]byte("allowance(address,address)"))[:10])
	t.Log(Keccak256Hash([]byte("balanceOf(address)"))[:10])

	t.Log(Keccak256Hash([]byte("transfer(address,uint256)"))[:10])
	t.Log(Keccak256Hash([]byte("transferFrom(address,address,uint256)"))[:10])
	t.Log(Keccak256Hash([]byte("approve(address,uint256)"))[:10])
}

func TestPubkeyToAddress(t *testing.T) {
	key, _ := HexToECDSA("7bbfec284ee43e328438d46ec803863c8e1367ab46072f7864c07e0a03ba61fd")
	t.Log(PubkeyToAddress(key.PubKey()))
	// 0x394586580ff4170c8a0244837202cbabe9070f66
	msg := "hello"
	hash := Keccak256([]byte(msg))
	sig, err := Sign(hash[:], key)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%x", sig)
	pub, err := SigToPub(hash, sig)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(PubkeyToAddress(pub))
}
