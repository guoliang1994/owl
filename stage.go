package owl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/dig"
	"net/http"
	_ "net/http/pprof"
	"os"
	"owl/log"
	"owl/utils/file"
	"path/filepath"
)

const (
	StoragePath   = "storage"
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
	runDir      string // 运行程序的目录
	binDir      string // 程序所在目录
	ConfManager *ConfManager
	AppCfg      GetConfigFunc
}

func New(appCfg string) *Stage {

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
		runDir:    runDir,
		binDir:    binDir,
	}

	_ = stage.Provide(func() *gin.Engine {
		e := gin.Default()
		e.Static(ResourcesPath, ResourcesPath)
		return e
	})

	_ = stage.Provide(NewLoggerFactory)
	configAbsPath := stage.RuntimePath("conf")
	file.CreateDirIfNotExists(configAbsPath)

	dataAbsPath := stage.RuntimePath("resource")
	file.CreateDirIfNotExists(dataAbsPath)

	cfgManager := NewConfigManager(configAbsPath)
	stage.ConfManager = cfgManager
	stage.AppCfg = cfgManager.GetConfig(appCfg).Get
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

// AbsBinDir 获取程序所在目录
func (i *Stage) AbsBinDir() string {
	return i.binDir
}

func (i *Stage) runPProf() {
	go func() {
		port := i.AppCfg("pprof-port").ToString()
		log.PrintLnBlue("pprof 监听端口：", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
			log.PrintLnRed("pprof 监听失败")
		}
	}()
}
