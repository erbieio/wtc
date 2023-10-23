package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Value represents a decoded RLP value, which is either a String or List.
type Value struct {
	// Only one of String or List is valid.  If String is "" then List is assumed valid.
	String string
	List   []Value
}

func (v *Value) IsList() bool {
	return v.String == ""
}

func (v *Value) IsString() bool {
	return v.String != ""
}

// Encode returns the 0x prefixed hex string of the RLP value
func (v Value) Encode() (string, error) {
	// If String is value encode that
	if v.IsString() {
		if !strings.HasPrefix(v.String, "0x") {
			return "", fmt.Errorf("invalid string value before encoding")
		}

		b, err := hex.DecodeString(v.String[2:])
		if err != nil {
			return "", fmt.Errorf("could not decode string value err: %v", err)
		}

		switch {
		case len(b) == 1 && b[0] <= 0x7f:
			// then the string is it's own encoding
			return v.String, nil
		case len(b) < 56:
			return "0x" + asUnprefixedHex(uint64(0x80+len(b))) + v.String[2:], nil
		default:
			size := asUnprefixedHex(uint64(len(b)))
			sizeSize := uint64(len(size) / 2)
			return "0x" + asUnprefixedHex(0xb7+sizeSize) + size + v.String[2:], nil
		}
	} else {
		// Otherwise encode the list, even if empty
		count := len(v.List)
		if count == 0 {
			// return the empty list
			return "0xc0", nil
		}

		data := make([]string, len(v.List))
		for i, item := range v.List {
			encoded, err := item.Encode()
			if err != nil {
				return "", fmt.Errorf("could not encode child item err: %v", err)
			}

			// Discard the 0x prefix
			data[i] = encoded[2:]
		}

		body := strings.Join(data, "")
		bodySize := uint64(len(body) / 2)
		if bodySize < 56 {
			// 0xc0 + bodySize
			return "0x" + asUnprefixedHex(bodySize+0xc0) + body, nil
		} else {
			bodySizeHex := asUnprefixedHex(bodySize)
			bodySizeSize := uint64(len(bodySizeHex) / 2)
			return "0x" + asUnprefixedHex(bodySizeSize+0xf7) + bodySizeHex + body, nil
		}
	}
}

// asUnprefixedHex converts a uint64 to a hex string WITHOUT the 0x prefix
func asUnprefixedHex(i uint64) string {
	bn := big.NewInt(0).SetUint64(i)
	b := bn.Bytes()
	h := hex.EncodeToString(b)
	return h
}

// HashToBytes returns the keccak256 hash of the encoded RLP as a byte slice
func (v Value) HashToBytes() ([]byte, error) {
	// TODO: Consider operating on the already encoded string vs. encoding inside this function
	// Encode the value back to a hex string
	var temp []byte
	encoded, err := v.Encode()
	if err != nil {
		return temp, fmt.Errorf("could not encode RLP value err: %v", err)
	}

	// Convert the string to bytes
	input := strings.Replace(strings.ToLower(encoded), "0x", "", 1)
	b, err := hex.DecodeString(input)
	if err != nil {
		return temp, fmt.Errorf("could not convert encoded to bytes err: %v", err)
	}

	// And feed the bytes into our hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write(b)
	sum := hash.Sum(nil)

	return sum, nil
}
