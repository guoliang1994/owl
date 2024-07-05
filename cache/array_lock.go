package cache

import "context"

type ArrayLock struct {
}

func (a ArrayLock) Get(ctx context.Context, callback func()) error {
	//TODO implement me
	panic("implement me")
}

func (a ArrayLock) Block(ctx context.Context, seconds int, callback func()) error {
	//TODO implement me
	panic("implement me")
}

func (a ArrayLock) Release() bool {
	//TODO implement me
	panic("implement me")
}

func (a ArrayLock) Owner() string {
	//TODO implement me
	panic("implement me")
}

func (a ArrayLock) ForceRelease() {
	//TODO implement me
	panic("implement me")
}
