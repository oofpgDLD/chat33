package model

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	var cfg types.Config
	if _, err := toml.DecodeFile("../etc/config.toml", &cfg); err != nil {
		panic(err)
	}
	db.InitDB(&cfg)
}

func Test_RevokeMsg(t *testing.T) {
	errMsg := RevokeMsg("8", "3141", 2)
	fmt.Println("err:", errMsg)
	fmt.Println("ok")
}

func Test_ReadSnapMsg(t *testing.T) {
	errMsg := ReadSnapMsg("8", "14797", types.IsRoom)
	fmt.Println("err:", errMsg)
	fmt.Println("ok")
}

func Test_ForwardMsg(t *testing.T) {
	info, errMsg := ForwardMsg("7", "363", types.IsRoom, 2, []string{"23843", "23845"}, []string{"421"}, []string{})
	fmt.Println("", info, errMsg)
	fmt.Println("ok")
}
func Test_ForwardEncryptMsg(t *testing.T) {
	//room测试数据
	room := make([]map[string]interface{}, 0)
	messages := make([]map[string]interface{}, 0)
	msg := make(map[string]interface{})
	v1 := make(map[string]interface{})
	v2 := make(map[string]interface{})
	msg["kid"] = "1567406864048"
	msg["encryptedMsg"] = "ee6e2dbd58b7c55e208f40c1fcba4ac1c32dee11a6aa199e67ffcc7d1adb506f70b1d959b66b4c19ad4106b2b9b03cca9b4a85943f5db85a15b750264f1689c63a2398bf6aa804eb48699bbb32bb21dc9262cc5639ac40efdd0b55536d05f23222667f10f7e6a2011d7e4cf3485c10941e191379c974c804a1a3856848be99d533829bf7976b278acda526fa42ce6b6b9e560aad3127ac2044dce99a3e74fe6376ed4a9791390129db76d45a1313f1881e3947097174cb439aa6"
	v1["msgType"] = 1
	v1["msg"] = msg
	messages = append(messages, v1)
	v2["targetId"] = "1047"
	v2["messages"] = messages
	room = append(room, v2)
	//user测试数据
	user := make([]map[string]interface{}, 0)
	usermessages := make([]map[string]interface{}, 0)
	usermsgs := make(map[string]interface{})
	u1 := make(map[string]interface{})
	u2 := make(map[string]interface{})
	usermsgs["fromKey"] = "033cb0196863cf872dde4d1a9e00ed0c6d7739be4e1a1a82233add8194ddfdb232"
	usermsgs["toKey"] = "02fd9d146dbcbdbcfe26ba61e182af2b5e4e9f6841703999de3f6365f3b3976ccf"
	usermsgs["encryptedMsg"] = "3bd4e50dbdb398823e4233dc2dd8f64292f22e7379db96fa3c28090ac9f78f96cbed7f23830b2524b7ac6f1e74a851dd6bb781485309921e00777c2635d036faa491b44681de888604e4ed7ec3e745e659992158588bac81873f3b57fddc75e6babb79ecc5ab31f26b379da7a4b5e19a1f76429e164af3fc679dc34dc02f968c195995bf263252105701b0de0480cd4473109aa52030bd6cad2e3198ea252d3774d704e986b43202f277b5776b273f09e2aa4c358a346cdfda4a"
	u1["msgType"] = 1
	u1["msg"] = usermsgs
	usermessages = append(usermessages, u1)
	u2["targetId"] = "8"
	u2["messages"] = usermessages
	user = append(user, u2)

	postParams := make(map[string]interface{})
	postParams["type"] = 1
	postParams["roomLogs"] = room
	postParams["userLogs"] = user
	info, errMsg := ForwardEncryptMsg("7", 1, room, user)
	fmt.Println("", info, errMsg)
	fmt.Println(postParams)
	fmt.Println("ok")
}

func Test_ClearlySearch(t *testing.T) {
	ret, errMsg := ClearlySearch("1004", "166", "469685076")
	fmt.Println(ret, errMsg)
	fmt.Println("ok")
}
