package routers

import (
	"github.com/astaxie/beego"
	"homework/backend/controllers"
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

	/*
		/order
	*/
	orderrepo := repositories.NewOrderManagerRepository("order", db)
	ordersevice := services.NewOrderService(orderrepo)
	ordercontroller := &controllers.OrderController{OrderService: ordersevice}
	beego.Router("/order/list", ordercontroller, "get:GetList")
	beego.Router("/order/deliver", ordercontroller, "post:GetUpdate")

	/*
		/product
	*/
	productrepo := repositories.NewProductManager("product", db, conn)
	productservice := services.NewProductService(productrepo)
	productcontroller := &controllers.ProductController{ProductService: productservice}
	beego.Router("/", productcontroller, "get:GetAll")
	beego.Router("/product/list", productcontroller, "get:GetAll")
	beego.Router("/product/update", productcontroller, "get:GetManager;post:PostManager")
	beego.Router("/product/insert", productcontroller, "get:GetInsert;post:PostInsert")
	beego.Router("/product/delete", productcontroller, "get:GetDelete")

	/*
		/shop
	*/
	shoprepo := repositories.NewShopRepository(db, "shop")
	shopservice := services.NewShopService(shoprepo)
	shopController := &controllers.ShopController{ShopService: shopservice}
	beego.Router("/shop/register", shopController, "get:GetRegister;post:PostRegister")
	beego.Router("/shop/login", shopController, "get:GetLogin;post:PostLogin")
	beego.Router("/shop/logout", shopController, "get:GetLogout")
}
