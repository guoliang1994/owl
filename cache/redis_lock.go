package cache

import (
	"context"
	"github.com/go-redis/redis"
)

// 实现 contract\cache\lock 接口
type RedisLock struct {
	c *redis.Client
}

func (r RedisLock) Get(ctx context.Context, callback func()) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisLock) Block(ctx context.Context, seconds int, callback func()) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisLock) Release() bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisLock) Owner() string {
	//TODO implement me
	panic("implement me")
}

func (r RedisLock) ForceRelease() {
	//TODO implement me
	panic("implement me")
}
