package common

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	. "strings"

	"github.com/fff-chain/3f-chain/core/common/hexutil"
)

// FFFAddress represents the 20 byte address of an Ethereum account.
type FFFAddress [AddressLength]byte

// BytesToAddress returns FFFAddress with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToFFFAddress(b []byte) FFFAddress {
	var a FFFAddress
	a.SetBytes(b)
	return a
}

func FFFAddressToAddress(address FFFAddress) Address {
	return BytesToAddress(address.Bytes())
}

func AddressToFFFAddress(address Address) FFFAddress {
	return BytesToFFFAddress(address.Bytes())
}

func AddressesToFFFAddresses(address []Address) []FFFAddress {
	var fffAddress []FFFAddress
	for z := range address {
		fffAddress = append(fffAddress, AddressToFFFAddress(address[z]))
	}
	return fffAddress
}

func FFFAddressesToAddresses(address []FFFAddress) []Address {
	var fffAddress []Address
	for z := range address {
		fffAddress = append(fffAddress, FFFAddressToAddress(address[z]))
	}
	return fffAddress
}

// BigToAddress returns FFFAddress with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToFFFAddress(b *big.Int) FFFAddress {
	return BytesToFFFAddress(b.Bytes())
}

// HexToAddress returns FFFAddress with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToFFFAddress(s string) FFFAddress {
	return BytesToFFFAddress(FFFFromHex(s))
}

func FFFFromHex(s string) []byte {
	return FromHex(FFFAddressDecode(s))
}

// IsHexFFFAddress verifies whether a string can represent a valid hex-encoded
// FFF address or not.
func IsHexFFFAddress(s string) bool {
	return IsHexAddress(FFFAddressDecode(s))
}

// Bytes gets the string representation of the underlying address.
func (a FFFAddress) Bytes() []byte { return a[:] }

// Hash converts an address to a hash by left-padding it with zeros.
func (a FFFAddress) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a FFFAddress) Hex() string {
	return FFFAddressEncode(BytesToAddress(a.Bytes()).Hex())
}

// String implements fmt.Stringer.
func (a FFFAddress) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter.
// FFFAddress supports the %v, %s, %v, %x, %X and %d format verbs.
func (a FFFAddress) Format(s fmt.State, c rune) {
	s.Write([]byte(a.String()))
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a), b will be cropped from the left.
func (a *FFFAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a FFFAddress) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalText parses a hash in hex syntax.
func (a *FFFAddress) UnmarshalText(input []byte) error {
	s := HexToAddress(FFFAddressDecode(string(input)))
	a.SetBytes(HexToAddress(FFFAddressDecode(string(input))).Bytes())
	return hexutil.UnmarshalFixedText("Address", []byte(s.hex()), a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *FFFAddress) UnmarshalJSON(input []byte) error {

	s := string(input)
	s = ReplaceAll(s, `"`, ``)
	return a.UnmarshalText([]byte(s))
}

// Scan implements Scanner for database/sql.
func (a *FFFAddress) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into FFFAddress", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into FFFAddress, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a FFFAddress) Value() (driver.Value, error) {
	return a[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (a FFFAddress) ImplementsGraphQLType(name string) bool { return name == "FFFAddress" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (a *FFFAddress) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = a.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for FFFAddress", input)
	}
	return err
}
