package helpers

import (
	"log"
	"net/url"
	"os"
	"path"
)

// AbsolutePath возвращает путь до файла в зависимости от режима запуска программы.
func AbsolutePath(pathStart string, pathEnd string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Server: can't get rooted path")
	}
	p := path.Base(cwd)

	if p == "go-pass_keeper" {
		return path.Join(pathStart, cwd, pathEnd)
	} else if p == "httplayer" || p == "appplayer" || p == "storelayer" {
		absPath, _ := url.JoinPath(pathStart, cwd, "../../..", pathEnd)
		return absPath
	} else {
		absPath, _ := url.JoinPath(pathStart, cwd, "../..", pathEnd)
		return absPath
	}
}
