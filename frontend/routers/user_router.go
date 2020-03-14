package routers

import (
	"github.com/astaxie/beego"
	"homework/common"
	"homework/frontend/controllers"
	"homework/models/repositories"
	"homework/models/services"
)

func init() {
	db, err := common.NewMysqlConn()
	if common.CheckErr(err) {
		return
	}
	userrepo := repositories.NewUserRepository("user", db)
	userservice := services.NewUserService(userrepo)
	usercontroller := &controllers.UserController{UserService: userservice}
	beego.Router("/user/register", usercontroller, "get:GetRegister;post:PostRegister")
	beego.Router("/user/login", usercontroller, "get:GetLogin;post:PostLogin")
}
