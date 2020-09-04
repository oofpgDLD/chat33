package address

import (
	"encoding/hex"
	"testing"
)

func Test_Encoding(t *testing.T) {
	in, err := hex.DecodeString("03f21818f989b977ed6b789ce088119d7ba5a4b37c39689ee2b4f273164555dbae")
	if err != nil {
		t.Error(err)
		return
	}
	addr := PublicKeyToAddress(NormalVer, in)
	t.Log(addr)
	if err := CheckAddress(NormalVer, addr); err != nil {
		t.Error(err)
		return
	}
	t.Log("check success")
}

func Test_CheckAddress(t *testing.T) {
	addr := "1JoFzozbxvst22c2K7MBYwQGjCaMZbC5Qm"
	if err := CheckAddress(NormalVer, addr); err != nil {
		t.Error(err)
		return
	}
	t.Log("check success")
}
