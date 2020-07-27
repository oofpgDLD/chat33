package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

var testpools []*redis.Pool

func init() {
	pool0 := &redis.Pool{
		MaxActive:   50,
		MaxIdle:     30,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://127.0.0.1:6379")
			if err != nil {
				panic(err)
			}
			//验证密码
			//if _,authErr := c.Do("AUTH",conf.Redis.Password);authErr != nil{
			//	panic(err)
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	testpools = append(testpools, pool0)
}

func Test_GetUserInfoById(t *testing.T) {
	userInfo, err := c.GetUserInfoById("1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userInfo)
}

func Test_UpdateAvatar(t *testing.T) {
	err := c.UpdateUserInfo("142", "username", "hahahah")
	fmt.Println(err)
}

func Test_Cmd(t *testing.T) {

	conn := testpools[0].Get()
	key := "test" + "123"
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Error("Close err", err)
		}
	}()

	/*_, err := conn.Do("HMSET", key,
		"test1", "nihao",
	)
	if err != nil {
		t.Error("HMSET EXPIRE err", err)
		return
	}*/

	err := conn.Send("MULTI")
	if err != nil {
		t.Error("send 'MULTI' err", err)
		return
	}

	/*_, err = conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
	if err != nil {
		t.Error("UpdateUserInfo EXPIRE err", err)
		return
	}
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		t.Error("UpdateUserInfo hexists err ", err)
		return
	}
	if count == 0 {
		return
	}*/
	_, err = conn.Do("hset", key, "test1", "hahahh")

	_, err = conn.Do("EXEC")
	if err != nil {
		t.Error("send 'EXEC' err", err)
		return
	}
	t.Log()
}
