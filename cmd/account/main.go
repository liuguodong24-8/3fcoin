package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/fff-chain/3f-chain/core/accounts/keystore"
	"github.com/fff-chain/3f-chain/core/common"
	"github.com/fff-chain/3f-chain/core/p2p/enode"
)

const password = "123456"

func main() {
	account, err := keystore.StoreKey("./key", password, 2, 1)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	keyjson, _ := ioutil.ReadFile("./" + account.URL.Path)
	pk, _ := keystore.DecryptKey(keyjson, password)
	enodeStr := enode.NewV4(&pk.PrivateKey.PublicKey, net.IP{127, 0, 0, 1}, 30300, 0)

	str := fmt.Sprintf("fff_addr =>:%s\neth_addr =>:%s\npassword =>:%s\npath     =>:%s\npk       =>:%s\nenode    =>:%s\n", pk.Address.Hex(), common.FFFAddressDecode(pk.Address.Hex()), password, account.URL.Path, hex.EncodeToString(pk.PrivateKey.D.Bytes()), enodeStr)

	ioutil.WriteFile("./key/"+pk.Address.Hex(), []byte(str), 0777)

	fmt.Println(str)
}
