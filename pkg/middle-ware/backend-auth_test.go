package middle_ware

import (
	"testing"
)

func Test_BcoinBackend(t *testing.T) {
	params := map[string]string{
		"rewardType": "2",
		"search":     "40",
		"size":       "2",
	}
	t.Log(BcoinBackend("wallet", "K2GMTloGU2ossdVV", "1562754897", params))
}
