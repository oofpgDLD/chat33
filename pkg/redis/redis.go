package redis

import (
	"github.com/garyburd/redigo/redis"
)

func NewRedisConn(conn redis.Conn) *RedisConn {
	return &RedisConn{
		c: conn,
	}
}

type RedisConn struct {
	c redis.Conn
}

func (conn *RedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	return conn.c.Do(commandName, args)
}

func (conn *RedisConn) Send(commandName string, args ...interface{}) error {
	return conn.c.Send(commandName, args)
}

func (conn *RedisConn) Close() error {
	return conn.c.Close()
}

func (conn *RedisConn) begin() error {
	return conn.c.Send("MULTI")
}

func (conn *RedisConn) NewTx() (*RedisTx, error) {
	err := conn.begin()
	if err != nil {
		return nil, err
	}
	return &RedisTx{
		RedisConn: conn,
	}, nil
}

type RedisTx struct {
	*RedisConn
}

func (t *RedisTx) RollBack() {
	defer func() {
		err := t.Close()
		if err != nil {

		}
	}()
	_, err := t.Do("DISCARD")
	if err != nil {

	}
}

func (t *RedisTx) Commit() error {
	defer func() {
		err := t.Close()
		if err != nil {

		}
	}()
	_, err := t.Do("EXEC")
	return err
}
