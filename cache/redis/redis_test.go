package redis

import (
	"encoding/json"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func Test_Conn(t *testing.T) {
	conn := GetConn()
	var nameMap = make(map[string]string)
	var nameMap1 = make(map[string]string)
	_, err := conn.Do("SET", "name", "wsj")
	if err != nil {
		t.Error(err)
	}
	nameMap["w"] = "sj"
	nameMap["w1"] = "sj"
	nameMapStr, _ := json.Marshal(nameMap)
	_, err = conn.Do("set", "nameMap", nameMapStr)
	if err != nil {
		t.Error(err)
	}
	v, _ := redis.Bytes(conn.Do("GET", "nameMap"))
	err = json.Unmarshal(v, &nameMap1)
	if err != nil {
		t.Error(err)
	}
	t.Log(nameMap1)
}
