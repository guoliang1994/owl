package cache

// Repository 定义缓存仓库接口
type Repository interface {
	// Retrieve an item from the cache and delete it.
	//
	// @param  string  $key
	// @param  mixed  $default
	// @return mixed
	Pull(key string, defaultVal any) any

	// Store an item in the cache.
	//
	// @param  string  $key
	// @param  mixed  $value
	// @param  \DateTimeInterface|\DateInterval|int|null  $ttl
	// @return bool
	Put(key string, value any, ttl any) bool

	// Store an item in the cache if the key does not exist.
	//
	// @param  string  $key
	// @param  mixed  $value
	// @param  \DateTimeInterface|\DateInterval|int|null  $ttl
	// @return bool
	Add(key string, value any, ttl any) bool

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
	Forever(key string, value any) bool

	// Get an item from the cache, or execute the given Closure and store the result.
	//
	// @param  string  $key
	// @param  \DateTimeInterface|\DateInterval|int|null  $ttl
	// @param  \Closure  $callback
	// @return mixed
	Remember(key string, ttl any, callback func() any) any

	// Get an item from the cache, or execute the given Closure and store the result forever.
	//
	// @param  string  $key
	// @param  \Closure  $callback
	// @return mixed
	Sear(key string, callback func() any) any

	// Get an item from the cache, or execute the given Closure and store the result forever.
	//
	// @param  string  $key
	// @param  \Closure  $callback
	// @return mixed
	RememberForever(key string, callback func() any) any

	// Remove an item from the cache.
	//
	// @param  string  $key
	// @return bool
	Forget(key string) bool

	// Get the cache store implementation.
	//
	// @return \Illuminate\Contracts\Cache\Store
	GetStore() Store
}
