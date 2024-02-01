package cache

// Factory 存储工厂，通过名称获取到具体的缓存者
type Factory interface {
	Store(name string)
}
