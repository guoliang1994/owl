package init_project

import (
	"fmt"
	"os"
)

func NewProject(name string) {
	createDir(name + "/models")
	createDir(name + "/routers")
	createDir(name + "/controllers")
	createDir(name + "/middlewares")
	createDir(name + "/storage")
	createDir(name + "/storage/app")
	createDir(name + "/storage/framework")
	createDir(name + "/storage/logs")
	createDir(name + "/conf")
}

func createDir(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
