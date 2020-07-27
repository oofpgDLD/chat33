package db

import "testing"

func Test_AdminCheckLogin(t *testing.T) {
	maps, err := AdminCheckLogin("1001", "admin")
	t.Log(maps)
	t.Log(err)
}
