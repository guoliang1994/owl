package middleware

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

const INDEX = "index.html"

var DistFs embed.FS

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}
type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}
func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

// FrontLocal  returns a middleware handler that serves static files in the given directory.
func FrontLocal(urlPrefix, root string) gin.HandlerFunc {
	return Front(urlPrefix, LocalFile(root, false))
}

// Front  returns a middleware handler that serves static files in the given directory.
func Front(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		} else {
			data, err := DistFs.ReadFile("dist/index.html")
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			c.Abort()
			return
		}
	}
}

type embedFileSystem struct {
	distFs embed.FS
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}
func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	DistFs = fsEmbed
	fsys, _ := fs.Sub(fsEmbed, targetPath)
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
