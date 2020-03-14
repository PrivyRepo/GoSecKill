package main

import (
	"github.com/astaxie/beego"
	_ "homework/frontend/routers"
)

func main() {
	beego.SetStaticPath("/public", "public")
	beego.Run()
}
