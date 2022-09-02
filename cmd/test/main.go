package main

import (
	"fmt"

	"github.com/liuguodong24-8/3fcoin/core/common"
)

func main() {
	newS := "\"FFF3QTZ3uQoVCiATg2ELuMjLb3SqoYtq6fnxV6jGMPFbLwJctj1q2qGj3F\""

	if common.IsHexAddress(newS[1 : len(newS)-1]) {
		fmt.Println(1)
	}

	fmt.Println(common.FFFAddressEncode("0x0d023dfc9c025e263d974985f3367d99f91e071b"))
	fmt.Println(common.FFFAddressDecode("FFF3QTZ3uQoVCiATg2ELuMjLb3SqoYtq6fnxV6jGMPFbLwJctj1q2qGj3F"))
	fmt.Println()
	// input = []byte(`"` + common.FFFAddressDecode(newS[1:len(newS)-1]) + `"`)

	// return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}
