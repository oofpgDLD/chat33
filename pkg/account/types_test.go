package account

import (
	"reflect"
	"testing"
)

func Test_Types(t *testing.T) {
	var u struct {
		BaseCurrency     string    `json:"baseCurrency" binding:"required"`
		BaseOpen         int       `json:"baseOpen"`
		RewardForUser    float64   `json:"rewardForUser" binding:"required"`
		RewardForInviter []float64 `json:"rewardForInviter" binding:"required"`
		AdvanceCurrency  string    `json:"advanceCurrency" binding:"required"`
		AdvanceOpen      int       `json:"advanceOpen"`
		ReachNum         int       `json:"reachNum"  binding:"required"`
		RewardForNum     struct {
			HHH string `json:"hhh" binding:"required"`
			SSS int    `json:"sss"`
		} `json:"rewardForNum"  binding:"required"`
	}
	/*	t := reflect.TypeOf(p)
		fmt.Println("Type: ", t.Name())
		v := reflect.ValueOf(t)
	*/
	u.BaseOpen = 1
	st := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		val := v.Field(i).Interface()
		t.Log(field.Tag.Get("json"), val)
	}
}

/*
func Test_Types2(t *testing.T) {
	var u struct {
		BaseCurrency     string    `json:"baseCurrency" binding:"required"`
		BaseOpen         int       `json:"baseOpen"`
		RewardForUser    float64   `json:"rewardForUser" binding:"required"`
		RewardForInviter []float64 `json:"rewardForInviter" binding:"required"`
		AdvanceCurrency  string    `json:"advanceCurrency" binding:"required"`
		AdvanceOpen      int       `json:"advanceOpen"`
		ReachNum         int       `json:"reachNum"  binding:"required"`
		RewardForNum     struct {
			HHH     	string    `json:"hhh" binding:"required"`
			SSS         int       `json:"sss"`
		}   `json:"rewardForNum"  binding:"required"`
	}

	ret := make(map[string]string)

	st := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	ss := st.Elem()
	vv := v.Elem()
	for i := 0; i < ss.NumField(); i++{
		field := ss.Field(i)
		val := vv.Field(i).Interface()

		ret[utility.ToString(field.Tag.Get("json"))] = utility.ToString(val)
	}

	return ret

}*/
