package file

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func CreateDirIfNotExists(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0777)
	}
}

// NormalizedPath 归一化处理路径，去除路径中错误的 \\ //
func NormalizedPath(path string) string {
	pattern := `[\\/]{1,}`
	re := regexp.MustCompile(pattern)        // 编译正则表达式
	result := re.ReplaceAllString(path, "/") // 使用正则表达式替换

	result = strings.TrimRight(result, "/") // 去除末尾 /
	return result
}

func DirIsEmpty(dirPath string) bool {

	// 获取文件夹中的文件列表
	fileList, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	// 检查文件列表是否为空
	if len(fileList) == 0 {
		return true
	} else {
		return false
	}
}
