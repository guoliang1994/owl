package owl

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/dig"
	_ "net/http/pprof"
	"os"
	"owl/utils/file"
	"path/filepath"
)

const (
	StoragePath   = "storage"
	ConfPath      = "conf"
	LogsPath      = StoragePath + "/logs"
	ResourcesPath = "resource"
	ViewsPath     = ResourcesPath + "/views"
)

type GetConfigFunc func(path ...interface{}) jsoniter.Any

type Dependency struct {
	Construct any
	Interface any
	Name      string
}

type Stage struct {
	*dig.Container
	runDir string // 运行程序的目录
	binDir string // 程序所在目录
}

func New() *Stage {

	var err error
	runDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	binDir := filepath.Dir(exePath)

	stage := &Stage{
		Container: dig.New(),
		runDir:    file.NormalizedPath(runDir),
		binDir:    file.NormalizedPath(binDir),
	}

	_ = stage.Provide(func() *gin.Engine {
		e := gin.Default()
		e.Static(ResourcesPath, stage.ResourcePath())
		return e
	})

	_ = stage.Provide(stage)
	_ = stage.Provide(NewLoggerFactory)
	_ = stage.Provide(NewConfigManager)

	configAbsPath := stage.RuntimePath("conf")
	file.CreateDirIfNotExists(configAbsPath)

	dataAbsPath := stage.RuntimePath("resource")
	file.CreateDirIfNotExists(dataAbsPath)

	return stage
}

// AbsRunDir 获取运行程序的目录
func (i *Stage) AbsRunDir() string {
	return i.runDir
}

// RuntimePath 获取配置所在目录
func (i *Stage) RuntimePath(relativePath string) (absPath string) {
	_, err := os.Stat(i.AbsBinDir() + "/" + relativePath)

	if os.IsNotExist(err) {
		return i.AbsRunDir() + "/" + relativePath
	} else {
		return i.AbsBinDir() + "/" + relativePath
	}
}

func (i *Stage) StoragePath() string {
	return i.RuntimePath(StoragePath)
}

func (i *Stage) ConfigPath() string {
	return i.RuntimePath(ConfPath)
}

func (i *Stage) ResourcePath() string {
	return i.RuntimePath(ResourcesPath)
}

func (i *Stage) LogPath() string {
	return i.RuntimePath(LogsPath)
}

// AbsBinDir 获取程序所在目录
func (i *Stage) AbsBinDir() string {
	return i.binDir
}

func (i *Stage) runPProf() {
	//go func() {
	//	port := i.AppCfg("pprof-port").ToString()
	//	log.PrintLnBlue("pprof 监听端口：", port)
	//	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
	//		log.PrintLnRed("pprof 监听失败")
	//	}
	//}()
}
