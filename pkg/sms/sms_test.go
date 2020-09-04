package sms

import (
	"testing"
)

/*
[smsConfig]
url="http://<host:port>/send/sms2"
codeType="chat_notice"
mobile=[""]
*/
func Test_Send(t *testing.T) {
	rlt, err := Send("http://118.31.52.32", "quick", "15763946517", "FzmRandom", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success", "rlt", rlt)
}

func Test_ValidateCode(t *testing.T) {
	err := ValidateCode("http://127.0.0.1/validate/code", "quick", "", "402805")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
