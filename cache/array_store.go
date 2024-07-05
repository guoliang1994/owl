package cache

type ArrayStore struct {
	store map[string][]byte
}

func (a ArrayStore) Get(key interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Many(keys []string) []interface{} {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Put(key string, value interface{}, seconds int) bool {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) PutMany(values []interface{}, seconds int) bool {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Increment(key string, value int) (int, bool) {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Decrement(key string, value int) (int, bool) {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Forever(key string, value interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Forget(key string) bool {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) Flush() bool {
	//TODO implement me
	panic("implement me")
}

func (a ArrayStore) GetPrefix() string {
	//TODO implement me
	panic("implement me")
}
