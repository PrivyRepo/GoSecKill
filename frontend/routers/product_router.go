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
	conn, err := common.NewRedisConn()
	if common.CheckErr(err) {
		return
	}
	productrepo := repositories.NewProductManager("product", db, conn)
	orderrepo := repositories.NewOrderManagerRepository("order", db)
	productservice := services.NewProductService(productrepo)
	orderservice := services.NewOrderService(orderrepo)
	productcontroller := &controllers.ProductController{ProductService: productservice, OrderService: orderservice}
	beego.Router("product/detail", productcontroller, "get:GetDetail")
	beego.Router("product/order", productcontroller, "get:GetOrder")
	beego.Router("product/list", productcontroller, "get:GetList")

}
