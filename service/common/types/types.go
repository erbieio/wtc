package types

import (
	"strconv"
)

// Data is a byte array represented by a prefixed hex string
type Data string

// UnmarshalJSON implements json.Unmarshaler.
func (d *Data) UnmarshalJSON(input []byte) error {
	return d.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Data) UnmarshalText(input []byte) error {
	*d = Data(input)
	return nil
}

type Address string

// UnmarshalJSON implements json.Unmarshaler.
func (a *Address) UnmarshalJSON(input []byte) error {
	return a.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler
func (a *Address) UnmarshalText(input []byte) error {
	*a = Address(input)
	return nil
}

type Hash string

// UnmarshalJSON implements json.Unmarshaler.
func (b *Hash) UnmarshalJSON(input []byte) error {
	return b.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *Hash) UnmarshalText(input []byte) error {
	*b = Hash(input)
	return nil
}

type Long int64

// UnmarshalJSON implements json.Unmarshaler.
func (b *Long) UnmarshalJSON(input []byte) error {
	if len(input) > 2 && input[0] == '"' {
		input = input[1 : len(input)-1]
	}
	return b.UnmarshalText(input)
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *Long) UnmarshalText(input []byte) error {
	value, err := strconv.ParseInt(string(input), 0, 64)
	*b = Long(value)
	return err
}
