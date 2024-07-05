package cache

// Store 定义缓存存储接口
type Store interface {
	// Retrieve an item from the cache by key.
	//
	// @param  string|array  $key
	// @return mixed
	Get(key interface{}) interface{}

	// Retrieve multiple items from the cache by key.
	// Items not found in the cache will have a null value.
	//
	// @param  array  $keys
	// @return array
	Many(keys []string) []interface{}

	// Store an item in the cache for a given number of seconds.
	//
	// @param  string  $key
	// @param  mixed  $value
	// @param  int  $seconds
	// @return bool
	Put(key string, value interface{}, seconds int) bool

	// Store multiple items in the cache for a given number of seconds.
	//
	// @param  array  $values
	// @param  int  $seconds
	// @return bool
	PutMany(values []interface{}, seconds int) bool

	// Increment the value of an item in the cache.
	//
	// @param  string  $key
	// @param  mixed  $value
	// @return int|bool
	Increment(key string, value int) (int, bool)

	// Decrement the value of an item in the cache.
	//
	// @param  string  $key
	// @param  mixed  $value
	// @return int|bool
	Decrement(key string, value int) (int, bool)

	// Store an item in the cache indefinitely.
	//
	// @param  string  $key
	// @param  mixed  $value
	// @return bool
	Forever(key string, value interface{}) bool

	// Remove an item from the cache.
	//
	// @param  string  $key
	// @return bool
	Forget(key string) bool

	// Remove all items from the cache.
	//
	// @return bool
	Flush() bool

	// Get the cache key prefix.
	//
	// @return string
	GetPrefix() string
}
