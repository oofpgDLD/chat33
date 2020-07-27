package redis

import (
	"hash/crc32"
	"sort"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var poolsIp []string //节点ip端口字符串
var pools []*redis.Pool
var nodeMap map[int]*redis.Pool //虚拟节点映射表
var vNodes []int                //虚拟节点

type nodeConfig struct {
	Url         string `json:"url"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
}

//初始化连接池
func InitRedis(conf *nodeConfig) {
	//加载配置
	pool0 := &redis.Pool{
		MaxActive:   conf.MaxActive,
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(conf.Url)
			if err != nil {
				panic(err)
			}
			//验证密码
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					panic(err)
				}
			}
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
	pools = append(pools, pool0)
	poolsIp = append(poolsIp, conf.Url)

	//一个节点对应的虚拟节点数
	var vnum = 10
	nodeMap = make(map[int]*redis.Pool)
	//加载节点映射表
	for k, v := range poolsIp {
		for i := 0; i < vnum; i++ {
			hash := GetHash(v + "#vn" + strconv.Itoa(i))
			vNodes = append(vNodes, hash)
			nodeMap[hash] = pools[k]
		}
	}

	//对虚拟节点进行排序
	sort.Ints(vNodes)
}

func GetConn() redis.Conn {
	return pools[0].Get()
}

func GetConnByKey(key string) redis.Conn {
	hash := GetHash(key)
	//查找虚拟节点
	vnode := binarySearch(vNodes, len(vNodes), hash)
	pool := nodeMap[vnode]
	return pool.Get()
}

//二分查第一个大于某值的数据
func binarySearch(nodes []int, n, v int) int {
	low := 0
	high := n - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nodes[mid] >= v {
			if mid == 0 || nodes[mid-1] < v {
				return nodes[mid]
			} else {
				high = mid - 1
			}
		} else {
			low = mid + 1
		}
	}
	return nodes[0]
}

func GetHash(str string) int {
	return int(crc32.ChecksumIEEE([]byte(str)))
}

func IsExists(reply interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		if len(reply) == 0 {
			return false, nil
		} else {
			return true, nil
		}
	case nil:
		return false, nil
	}
	return true, nil
}
