package routers

import (
	"github.com/astaxie/beego"
	"homework/backend/controllers"
	"homework/common"
	"homework/models/repositories"
	"homework/models/services"
)

func init() {
	db, err := common.NewMysqlConn()
	if common.CheckErr(err) {
		return
	}
	shoprepo := repositories.NewShopRepository(db, "shop")
	shopservice := services.NewShopService(shoprepo)
	shopController := &controllers.ShopController{ShopService: shopservice}
	beego.Router("/shop/register", shopController, "get:GetRegister;post:PostRegister")
	beego.Router("/shop/login", shopController, "get:GetLogin;post:PostLogin")
	beego.Router("/shop/logout", shopController, "get:GetLogout")
}
