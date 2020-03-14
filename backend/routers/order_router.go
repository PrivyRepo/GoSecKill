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
	orderrepo := repositories.NewOrderManagerRepository("order", db)
	ordersevice := services.NewOrderService(orderrepo)
	ordercontroller := &controllers.OrderController{OrderService: ordersevice}
	beego.Router("/order/list", ordercontroller, "get:GetList")
	//beego.Router("/order/update",ordercontroller,"post:UpdateOrder")
	beego.Router("/order/delete", ordercontroller, "get:GetDlete")

}
