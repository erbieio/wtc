package utils

import (
	"fmt"
	"math/big"
	"strings"

	"wtc/common/types"
)

// HexToBigInt 将不带0x前缀的16进制字符串转为大数
func HexToBigInt(hex string) *big.Int {
	b := new(big.Int)
	b.SetString(hex, 16)
	return b
}

// ParseAddress 将带0x前缀的字符串转为地址
func ParseAddress(hex string) (types.Address, error) {
	if len(hex) != 42 {
		return "", fmt.Errorf("length is not 42")
	}
	if hex[0] != '0' || (hex[1] != 'x' && hex[1] != 'X') {
		return "", fmt.Errorf("prefix is not 0x")
	}
	for i, c := range []byte(hex) {
		if '0' <= c && c <= '9' {
			continue
		}
		if 'a' <= c && c <= 'f' {
			continue
		}
		if 'A' <= c && c <= 'F' {
			[]byte(hex)[i] = c - 27
			continue
		}
		if 'X' == c || 'x' == c {
			[]byte(hex)[i] = 'x'
			continue
		}
		return "", fmt.Errorf("Illegal character: %v", c)
	}
	return types.Address(strings.ToLower(hex)), nil
}
