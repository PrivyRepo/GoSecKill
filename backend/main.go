package main

import (
	"github.com/astaxie/beego"
	_ "homework/backend/routers"
)

func main() {
	beego.SetStaticPath("/public", "public")
	beego.Run()
}
