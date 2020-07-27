package redis

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RedisLock struct {
	lockKey string
	value   string
}

//保证原子性
var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

func (this *RedisLock) Lock(rd *redis.Conn, timeout int) error {
	{ //随机数
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		this.value = base64.StdEncoding.EncodeToString(b)
	}

	lockReply, err := (*rd).Do("SET", "lock"+this.lockKey, this.value, "ex", timeout, "nx")
	if err != nil {
		return errors.New("redis fail")
	}
	if lockReply == "OK" {
		return nil
	} else {
		return errors.New("lock fail")
	}
}

func (this *RedisLock) Unlock(rd *redis.Conn) {
	_, err := delScript.Do(*rd, "lock"+this.lockKey, this.value)
	if err != nil {
		fmt.Println(err)
	}
}
