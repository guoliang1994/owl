package contract

type ServiceProvider interface {
	Register() // 构建容器需要的对象
	Boot()
}
