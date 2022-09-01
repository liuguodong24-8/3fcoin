package main

import (
	"fmt"

	"github.com/fff-chain/3f-chain/core/common"
)

func main() {
	UnmarshalFixedUnprefixedText()
}

func UnmarshalFixedUnprefixedText() {
	tests := []struct {
		input   string
		want    []byte
		wantErr error
	}{
		{input: "0x3e09c89e643cc3f30601fa08adb07e87546e65f9", wantErr: common.ErrOddLength},
	}

	for _, test := range tests {
		out := make([]byte, 20)
		err := common.UnmarshalFixedText("common.Address", []byte(test.input), out)
		fmt.Println(err)
	}

}
