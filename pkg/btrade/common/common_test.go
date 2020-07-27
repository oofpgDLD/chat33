package common

import (
	"testing"
)

func TestStringToFloat(t *testing.T) {
	price := "9782.13"
	fprice := StringToFloat64(price)
	pp := int64(fprice * 1e8)
	t.Log(fprice, pp)
}

func TestString2BaseInt64(t *testing.T) {
	nums := []string{
		"0",
		"1",
		"98765423",
		"12",
		"123.4",
		"0.33",
		"0.00000041",
		"9782.12345678",
		"9782.13",
		"2.000059",
		"1.19",
	}
	for _, v := range nums {
		fprice := String2BaseInt64(v)
		t.Log(v, fprice)
	}
}
