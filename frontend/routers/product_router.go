package routers

import (
	"github.com/astaxie/beego"
	"homework/common"
	"homework/frontend/controllers"
	"homework/models/repositories"
	"homework/models/services"
	"homework/rabbitmq"
)

func init() {
	db, err := common.NewMysqlConn()
	if common.CheckErr(err) {
		return
	}
	conn := common.NewRedisConn()
	productrepo := repositories.NewProductManager("product", db, conn)
	orderrepo := repositories.NewOrderManagerRepository("order", db)
	productservice := services.NewProductService(productrepo)
	orderservice := services.NewOrderService(orderrepo)
	rabbitmq := rabbitmq.NewRabbitMQSimple("imooc")
	productcontroller := &controllers.ProductController{ProductService: productservice, OrderService: orderservice, RabbitMQ: rabbitmq}
	beego.Router("product/detail", productcontroller, "get:GetDetail")
	beego.Router("product/order", productcontroller, "get:GetOrder")
	beego.Router("product/list", productcontroller, "get:GetList")

	go rabbitmq.ConsumeSimple(orderservice, productservice)
}
