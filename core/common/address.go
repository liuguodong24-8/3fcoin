package common

import (
	"strings"
)

var (
	FFFHeader = "FFF"
	ETHHeader = "0x"
)

func FFFAddressEncode(hex string) string {
	hex = strings.ToLower(hex)
	var relHex = ""
	if strings.ToLower(hex[0:2]) == "0x" {
		relHex = hex[2:]
	} else {
		relHex = hex
	}
	return FFFHeader + Base58Encoding(relHex)
}

func FFFAddressDecode(hex string) string {
	var relHex = ""
	if strings.ToLower(hex[0:3]) != "fff" {
		relHex = hex
	} else {
		relHex = hex[3:]
	}
	return ETHHeader + Base58Decoding(relHex)
}
