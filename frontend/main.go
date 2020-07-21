package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "homework/frontend/middleware"
	_ "homework/frontend/routers"
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
	beego.SetViewsPath(GetAPPRootPath() + "/views/")
	beego.SetStaticPath("/public", GetAPPRootPath()+"/public")
	beego.Run()
}
