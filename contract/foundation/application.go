package foundation

// Application 定义应用接口
type Application interface {
	// Version 获取应用版本号
	Version() string

	// BasePath 获取 Laravel 安装的基础路径
	BasePath(path string) string

	// BootstrapPath 获取引导目录路径
	BootstrapPath(path string) string

	// ConfigPath 获取应用配置文件路径
	ConfigPath(path string) string

	// DatabasePath 获取数据库目录路径
	DatabasePath(path string) string

	// ResourcePath 获取资源目录路径
	ResourcePath(path string) string

	// StoragePath 获取存储目录路径
	StoragePath(path string) string

	// Environment 获取或检查当前应用环境
	Environment(environments ...string) (string, bool)

	// RunningInConsole 判断应用是否在控制台运行
	RunningInConsole() bool

	// RunningUnitTests 判断应用是否正在运行单元测试
	RunningUnitTests() bool

	// MaintenanceMode 获取维护模式管理器实例
	MaintenanceMode() MaintenanceMode

	// IsDownForMaintenance 判断应用是否处于维护状态
	IsDownForMaintenance() bool

	// RegisterConfiguredProviders 注册所有已配置的提供者
	RegisterConfiguredProviders()

	// Register 注册服务提供者
	Register(provider string, force bool) ServiceProvider

	// RegisterDeferredProvider 注册延迟的提供者和服务
	RegisterDeferredProvider(provider, service string)

	// ResolveProvider 根据名称解析服务提供者实例
	ResolveProvider(provider string) ServiceProvider

	// Boot 启动应用的服务提供者
	Boot()

	// Booting 注册新的启动监听器
	Booting(callback func())

	// Booted 注册新的已启动监听器
	Booted(callback func())

	// BootstrapWith 运行给定的引导类数组
	BootstrapWith(bootstrappers []string)

	// GetLocale 获取当前应用的区域设置
	GetLocale() string

	// GetNamespace 获取应用的命名空间
	GetNamespace() (string, error)

	// GetProviders 获取已注册的服务提供者实例
	GetProviders(provider string) []ServiceProvider

	// HasBeenBootstrapped 判断应用是否已引导
	HasBeenBootstrapped() bool

	// LoadDeferredProviders 加载并启动所有剩余的延迟提供者
	LoadDeferredProviders()

	// SetLocale 设置当前应用的区域设置
	SetLocale(locale string)

	// ShouldSkipMiddleware 判断是否禁用中间件
	ShouldSkipMiddleware() bool

	// Terminating 注册终止回调
	Terminating(callback interface{}) Application

	// Terminate 终止应用
	Terminate()
}

// MaintenanceMode 维护模式接口
type MaintenanceMode interface {
	// 维护模式相关的方法
}

// ServiceProvider 服务提供者接口
type ServiceProvider interface {
	// 服务提供者相关的方法
}
