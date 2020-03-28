package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "homework/frontend/server/middleware"
	_ "homework/frontend/server/routers"
	"os"
	"os/exec"
	"path/filepath"
)

func GetAPPRootPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}

func main() {
	logs.SetLogger("file", `{"filename":"error.log","level":5}`)
	beego.Run()
}
