package routers

import (
	"homework/frontend/controllers"
	"github.com/astaxie/beego"
	"homework/common/mysql"
	"homework/common/redis"
	"homework/models/repositories"
	"homework/models/services"
)

func init() {
	db, err := mysql.NewMysqlConn()
	if err != nil {
		return
	}
	conn := redis.NewRedisConn()
	productrepo := repositories.NewProductManager("product", db, conn)
	productservice := services.NewProductService(productrepo)
	//rabbitmq := rabbitmq2.NewRabbitMQSimple("imooc")
	productcontroller := &controllers.ProductController{ProductService: productservice}
	seckillcontroller := &controllers.SeckillController{MySqlConn: db}
	beego.Router("/product/detail", productcontroller, "get:GetTestDetail")
	beego.Router("/product/testdetail", productcontroller, "get:GetTestDetail")
	beego.Router("/product/list", productcontroller, "get:GetList")
	beego.Router("/product/kill", seckillcontroller, "get:Kill")
	//go rabbitmq.ConsumeSimple(orderservice, productservice)

	redisconn := redis.NewRedisConn()
	userrepo := repositories.NewUserRepository("user", db, redisconn)
	userservice := services.NewUserService(userrepo)
	usercontroller := &controllers.UserController{UserService: userservice}
	beego.Router("/user/register", usercontroller, "post:PostRegister;get:GetRegister")
	beego.Router("/user/login", usercontroller, "post:PostLogin;get:GetLogin")
	beego.Router("/user/logout", usercontroller, "get:GetLogout")

	//SecKillcontroller := &controllers.SeckillController{RabbitMQ: rabbitmq, MySqlConn: db}
	//beego.Router("/order/execution",SecKillcontroller,"get:Kill")

	orderrepo := repositories.NewOrderManagerRepository("order", db)
	orderservice := services.NewOrderService(orderrepo)
	orderController := &controllers.OrderController{OrderService: orderservice}
	beego.Router("/order/list", orderController, "get:OrderList")
	beego.Router("/order/payoff", orderController, "get:GetPayoff")
	beego.Router("/order/cancel", orderController, "get:GetCancel")
}
