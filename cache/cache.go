package cache

import "fmt"

var (
	drivers = make(map[string]cache)
)

var Cache cache

type cache interface {
	UserCacheI
	FriendCacheI
	RoomCacheI
	ApplyCacheI
	OrderCacheI
	AccountCacheI
	PariseCacheI
}

func Register(name string, driver cache) {
	if driver == nil {
		panic("cahce: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("cache: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func GetInstance(name string) (c cache, err error) {
	c, ok := drivers[name]
	if !ok {
		err = fmt.Errorf("unknown driver %q", name)
		return
	}
	return c, nil
}
