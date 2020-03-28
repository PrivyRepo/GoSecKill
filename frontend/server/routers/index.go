package routers

import (
	"github.com/astaxie/beego"
	"homework/common/mysql"
	rabbitmq2 "homework/common/rabbitmq"
	"homework/common/redis"
	"homework/frontend/server/controllers"
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
	rabbitmq := rabbitmq2.NewRabbitMQSimple("imooc")
	productcontroller := &controllers.ProductController{ProductService: productservice}
	beego.Router("/api/product/detail", productcontroller, "get:GetDetail")
	beego.Router("/api/product/list", productcontroller, "get:GetList")
	//go rabbitmq.ConsumeSimple(orderservice, productservice)

	redisconn := redis.NewRedisConn()
	userrepo := repositories.NewUserRepository("user", db, redisconn)
	userservice := services.NewUserService(userrepo)
	usercontroller := &controllers.UserController{UserService: userservice}
	beego.Router("/api/user/register", usercontroller, "post:PostRegister")
	beego.Router("/api/user/login", usercontroller, "post:PostLogin")
	beego.Router("/api/user/logout", usercontroller, "get:GetLogout")

	//SecKillcontroller := &controllers.SeckillController{RabbitMQ: rabbitmq, MySqlConn: db}
	//beego.Router("/order/execution",SecKillcontroller,"get:Kill")

	orderrepo := repositories.NewOrderManagerRepository("order", db)
	orderservice := services.NewOrderService(orderrepo)
	orderController := &controllers.OrderController{OrderService: orderservice}
	beego.Router("/order/list", orderController, "get:OrderList")
}
