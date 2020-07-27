package model

import (
	"testing"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func Test_Balance(t *testing.T) {
	ret, err := Balance("1006", "b2d71f69fc448b229a0f374df728d096e9a85e6c")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}

func Test_getRPH5Url(t *testing.T) {
	ret := getRPH5Url("1001", "1", "294f1930-e401-11e8-b91b-a9f3467a608c")
	t.Log(ret)
}

func Test_RedEnvelopeDetail(t *testing.T) {
	ret, err := RedEnvelopeDetail("1006", "385", "7c6118fc8469abb03e03d21cfcbd06af005847e2", "bde361fb-d9b1-4803-8918-48a337415c58")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}

func Test_Length(t *testing.T) {
	remark := "呵呵啥收拾收拾U盾好的好的好的好的几点还打不打黄赌毒度好的好的好的记得记得几点度好的好的好的好的12312312312312312312312的好的好的记得记得几点度好的好的好的好的12312312312312312312312"
	req := &types.SendParams{
		Remark: remark,
	}

	arry := []rune(utility.ToString(req.Remark))
	if len(arry) > types.RemarkLengthLimit {
		req.Remark = string(arry[:types.RemarkLengthLimit]) + "..."
	}
	t.Log(req.Remark)
}
