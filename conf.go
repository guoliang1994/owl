package owl

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"owl/log"
	"path/filepath"
	"strings"
)

type ConfManager struct {
	confDir      string
	changeNotify map[string]chan string
	allCfg       map[string]map[string]any // 存储所有的配置
}

var (
	CfgChangeNotify = make(map[string]chan string, 10) // 配置修改时通知
)

func NewConfigManager(confDir string) *ConfManager {
	manager := ConfManager{
		changeNotify: nil,
		allCfg:       make(map[string]map[string]any),
		confDir:      confDir,
	}
	err := filepath.Walk(confDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		cfgMap := make(map[string]any)
		ext := strings.Replace(filepath.Ext(info.Name()), ".", "", -1)
		name := strings.Replace(info.Name(), "."+ext, "", -1)
		absPath, v := manager.LoadConfig(name, ext, &cfgMap)
		cfgMap["abs-path"] = absPath
		cfgMap["vip"] = v
		manager.allCfg[name] = cfgMap
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &manager
}

func (i *ConfManager) GetConfig(key string) jsoniter.Any {
	var getter jsoniter.Any

	marshal, err := jsoniter.Marshal(i.allCfg)
	if err != nil {
		return getter
	}

	pathArr := strings.Split(key, ".")

	for _, path := range pathArr {
		if getter != nil {
			getter = jsoniter.Get([]byte(getter.ToString()), path)
		} else {
			getter = jsoniter.Get(marshal, path)
		}
	}
	return getter
}

func (i *ConfManager) SaveConfig(fileName string, key string, value any) {
	cfg, ok := i.allCfg[fileName]
	if ok {
		v := cfg["vip"].(*viper.Viper)
		v.Set(key, value)
		err := v.WriteConfig()
		if err != nil {
			fmt.Println("保存配置失败")
		}
	}
}

// LoadConfig 读取文件中的配置
func (i *ConfManager) LoadConfig(fileName, cfgType string, c any) (string, *viper.Viper) {

	confFilePath := fmt.Sprintf("%s/%s.%s", i.confDir, fileName, cfgType)

	v := viper.New()
	v.SetConfigType(cfgType)
	v.SetConfigFile(confFilePath)
	log.PrintLnBlue("配置文件: ", confFilePath)
	cfg, err := os.ReadFile(confFilePath)
	if err != nil {
		log.PrintRed("配置文件不存在")
		os.Exit(200)
	}
	v.AddConfigPath(confFilePath)
	if err = v.ReadConfig(bytes.NewReader(cfg)); err != nil {
		panic(fmt.Sprint("读取配置文件失败", err))
	}

	// 转换为结构体
	if err := v.Unmarshal(&c); err != nil {
		panic("转为配置结构体失败")
	}

	// Watch for changes in the config file
	v.WatchConfig()
	ch := make(chan string, 10)
	CfgChangeNotify[confFilePath] = ch
	// Register a callback function to handle the changes
	v.OnConfigChange(func(e fsnotify.Event) {
		CfgChangeNotify[confFilePath] <- confFilePath
	})
	return confFilePath, v
}
