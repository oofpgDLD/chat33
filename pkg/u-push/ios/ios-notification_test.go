package ios

import (
	"encoding/json"
	"testing"

	push "github.com/33cn/chat33/pkg/u-push"
)

func Test_IOSExtraReal(t *testing.T) {
	var client push.PushClient
	//unicast := NewAndroidUnicast("5d5230610cafb2025f00076e", "x853sko9sp7hbcrmzlhnnbyvqxx0nsni")
	unicast := NewIOSUnicast("5dc924cb570df3de690009d8", "tk4mbj5n4efcago13yp5if6tql57pkkx")

	//unicast.SetDeviceToken("Ahm4teWYFAA2cWZKozUSQj4nZcxg1YvAvQYNfJNFr1Gi")
	//unicast.SetDeviceToken("Ak9Xqj2AyYdwAmZeBP0MO-WxG_06Eawya-AtCpJ-9N6W")
	unicast.SetDeviceToken("a2cf98195d98233d11a17fd7e71c87f0e099c50546fc12d199ee84ebf04ea534")

	unicast.SetAlertJson(push.IOSAlert{
		Title: "中文的title",
		Body:  "iOS unicast text",
	})
	//unicast.SetBadge(0)
	unicast.SetSound("default")
	// 线上模式
	unicast.SetTestMode()
	// Set customized fields
	unicast.SetCustomizedField("targetId", "22")
	unicast.SetCustomizedField("channelType", "2")
	b, err := json.Marshal(unicast)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
	err = client.Send(unicast)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
