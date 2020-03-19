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
	conn := common.NewRedisConn()
	productrepo := repositories.NewProductManager("product", db, conn)
	productservice := services.NewProductService(productrepo)

	productcontroller := &controllers.ProductController{ProductService: productservice}
	beego.Router("/product/list", productcontroller, "get:GetAll")
	beego.Router("/product/update", productcontroller, "get:GetManager;post:PostManager")
	beego.Router("/product/insert", productcontroller, "get:GetInsert;post:PostInsert")
	beego.Router("/product/delete", productcontroller, "get:GetDelete")
}
