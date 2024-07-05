package cache

// LockProvider 定义锁提供器接口
type LockProvider interface {
	// Get a lock instance.
	//
	// @param  string  $name
	// @param  int  $seconds
	// @param  string|null  $owner
	// @return \Illuminate\Contracts\Cache\Lock
	Lock(name string, seconds int, owner string) Lock

	// Restore a lock instance using the owner identifier.
	//
	// @param  string  $name
	// @param  string  $owner
	// @return \Illuminate\Contracts\Cache\Lock
	RestoreLock(name, owner string) Lock
}
