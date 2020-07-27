package android

import (
	"encoding/json"
	"testing"

	push "github.com/33cn/chat33/pkg/u-push"
)

func Test_Android(t *testing.T) {
	appkey := "5d5ce767570df392c3000428"
	appMasterSecret := "cxeguyg71bdjxjxed0euzxwkyruueviw"

	var client push.PushClient
	unicast := NewAndroidUnicast(appkey, appMasterSecret)
	// TODO Set your device token
	unicast.SetDeviceToken("Ahm4teWYFAA2cWZKozUSQj7XRvxhXbWlia4xTrbY2Cc7")
	unicast.SetTicker("Android unicast ticker")
	unicast.SetTitle("中文的title")
	unicast.SetText("Android unicast text")
	unicast.GoAppAfterOpen()
	unicast.SetDisplayType(push.NOTIFICATION)
	// TODO Set 'production_mode' to 'false' if it's a test device.
	// For how to register a test device, please see the developer doc.
	unicast.SetReleaseMode()
	// Set customized fields
	unicast.SetExtraField("test", "helloworld")
	err := client.Send(unicast)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_AndroidExtra(t *testing.T) {
	appkey := "5d5230610cafb2025f00076e"
	appMasterSecret := "x853sko9sp7hbcrmzlhnnbyvqxx0nsni"

	var client push.PushClient
	unicast := NewAndroidUnicast(appkey, appMasterSecret)
	// TODO Set your device token
	unicast.SetDeviceToken("Ahm4teWYFAA2cWZKozUSQj4nZcxg1YvAvQYNfJNFr1Gi")
	unicast.SetTicker("Android unicast ticker")
	unicast.SetTitle("title: extra test")
	unicast.SetText("Android unicast extra text")
	//unicast.GoAppAfterOpen()
	unicast.GoCustomAfterOpen("")
	unicast.SetDisplayType(push.NOTIFICATION)
	// TODO Set 'production_mode' to 'false' if it's a test device.
	// For how to register a test device, please see the developer doc.
	unicast.SetReleaseMode()
	// Set customized fields
	unicast.SetExtraField("targetId", "49")
	unicast.SetExtraField("channelType", "2")
	err := client.Send(unicast)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_AndroidExtraReal(t *testing.T) {
	var client push.PushClient
	//unicast := NewAndroidUnicast("5d5230610cafb2025f00076e", "x853sko9sp7hbcrmzlhnnbyvqxx0nsni")
	unicast := NewAndroidUnicast("5d5ce767570df392c3000428", "cxeguyg71bdjxjxed0euzxwkyruueviw")

	//unicast.SetDeviceToken("Ahm4teWYFAA2cWZKozUSQj4nZcxg1YvAvQYNfJNFr1Gi")
	//unicast.SetDeviceToken("Ak9Xqj2AyYdwAmZeBP0MO-WxG_06Eawya-AtCpJ-9N6W")
	unicast.SetDeviceToken("Ahm4teWYFAA2cWZKozUSQj7XRvxhXbWlia4xTrbY2Cc7")
	unicast.SetTicker("Android unicast ticker")
	unicast.SetTitle("你的素颜如水群")
	unicast.SetText("chatjzUzoqQRWz: ...")
	unicast.GoCustomAfterOpen("")
	unicast.SetDisplayType(push.NOTIFICATION)
	unicast.SetMipush(true, "com.fzm.push.SystemPushActivity")
	// 线上模式
	unicast.SetReleaseMode()
	// Set customized fields
	unicast.SetExtraField("targetId", "22")
	unicast.SetExtraField("channelType", "2")
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
