package cache

import "github.com/go-redis/redis"

type RedisStore struct {
	client *redis.Client
}

func (r RedisStore) Get(key interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Many(keys []string) []interface{} {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Put(key string, value interface{}, seconds int) bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) PutMany(values []interface{}, seconds int) bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Increment(key string, value int) (int, bool) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Decrement(key string, value int) (int, bool) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Forever(key string, value interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Forget(key string) bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) Flush() bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisStore) GetPrefix() string {
	//TODO implement me
	panic("implement me")
}
