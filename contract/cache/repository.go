package cache

import (
	"context"
	"time"
)

// Store 类似于 PHP 中的 Illuminate\Contracts\Cache\Store 接口
type Store interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}, ttl time.Duration) error
	Add(key string, value interface{}, ttl time.Duration) error
	Increment(key string, value int64) (int64, error)
	Decrement(key string, value int64) (int64, error)
	Forever(key string, value interface{}) error
	Forget(key string) error
}

// Repository 是一个扩展自 Store 的缓存操作接口
type Repository interface {
	Store

	// Retrieve an item from the cache and delete it.
	Pull(key string, def interface{}) (interface{}, error)

	// Store an item in the cache if the key does not exist.
	AddOrSet(key string, value interface{}, ttl time.Duration) (bool, error)

	// Increment the value of an item in the cache.
	IncrementByKey(key string, value int64) (int64, error)

	// Decrement the value of an item in the cache.
	DecrementByKey(key string, value int64) (int64, error)

	// Store an item in the cache indefinitely.
	RememberForever(key string, callback func() (interface{}, error)) (interface{}, error)

	// Get an item from the cache, or execute the given function and store the result.
	Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)

	// Sear is likely a typo in the original PHP code. Assuming it should be RememberForever.
	RememberForeverWithContext(ctx context.Context, key string, callback func() (interface{}, error)) (interface{}, error)
}

// 示例方法实现（具体实现取决于实际存储和同步机制）
type MyRepository struct {
	store Store
}

func (r *MyRepository) Pull(key string, def interface{}) (interface{}, error) {
	value, err := r.store.Get(key)
	if err != nil {
		return def, nil
	}
	// r.Forget(key)
	return value, nil
}

// AddOrSet 方法合并了 add 和 put 功能
func (r *MyRepository) AddOrSet(key string, value interface{}, ttl time.Duration) (bool, error) {
	//_, existsErr := r.store.Get(key)
	//if existsErr == nil {
	//	return false, r.Put(key, value, ttl)
	//}
	//return true, r.Add(key, value, ttl)
	return true, nil
}

// 其他示例方法实现略...
