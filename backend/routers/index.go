package routers

import (
	"github.com/astaxie/beego"
	"homework/backend/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
}
