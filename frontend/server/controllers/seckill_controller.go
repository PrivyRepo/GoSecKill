package controllers

import (
	"database/sql"
	"github.com/astaxie/beego"
	"homework/common/rabbitmq"
	"homework/models/datamodels"
	"strconv"
)

type SeckillController struct {
	beego.Controller
	RabbitMQ  *rabbitmq.RabbitMQ
	MySqlConn *sql.DB
}

func (this *SeckillController) Kill() {

	productstring := this.GetString("productID")
	uidstring := this.Ctx.GetCookie("uid")
	productid, e := strconv.Atoi(productstring)
	if e != nil {
		this.Abort("500")
	}
	userid, _ := strconv.Atoi(uidstring)

	////创建消息体
	//message := datamodels.NewMessage(int64(userID), int64(productID))
	//bytes, e := json.Marshal(message)
	//if common.CheckErr(e) {
	//	this.Abort("401")
	//}
	//
	//e = this.RabbitMQ.PublishSimple(string(bytes))
	//if common.CheckErr(e) {
	//	this.Abort("401")
	//}
	//this.Ctx.WriteString("successful")

	//不使用事务保证数据一致性
	//	product, e := this.ProductService.GetProductByID(int64(productid))
	//	if e!=nil  {
	//		this.Abort("500")
	//	}
	//	var showMessage string = "抢购失败"
	//	//判断商品数量是否满足需求
	//	if product.ProductNum > 0 {
	//		//扣除商品数量
	//		product.ProductNum -= 1
	//		err := this.ProductService.UpdateProduct(product)
	//		if err!=nil {
	//			this.Abort("501")
	//		}
	//		//创建订单
	//		order := &datamodels.Order{
	//			UserId:      int64(uid),
	//			ProductId:   int64(productid),
	//			OrderPayStatus: datamodels.PayWait,
	//		}
	//		_, err = this.OrderService.InsertOrder(order)
	//		if err!=nil {
	//			return
	//		} else {
	//			showMessage = "抢购成功"
	//		}
	//	}

	//使用事务保证数据一致性
	result := this.MySqlConn.QueryRow("SELECT productNum FROM `product` WHERE ID = ?", productstring)
	if e != nil {
		this.Abort("501")
	}
	var productNum int
	e = result.Scan(&productNum)
	if e != nil {
		this.Abort("501")
	}
	var showMessage string = "抢购失败"
	if productNum > 0 {
		tx, e := this.MySqlConn.Begin()
		if e != nil {
			this.Abort("500")
		}
		updateSQL := "UPDATE `product` SET productNum = productNum -1 WHERE ID = ?"
		_, e2 := tx.Exec(updateSQL, productstring)
		insertSQL := "INSERT `order` set userID=?,productID=?,orderPayStatus=?,orderDeliverStatus=?"
		_, e1 := tx.Exec(insertSQL, userid, productid, datamodels.PayWait, datamodels.DeliverWait)
		if e1 != nil || e2 != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			showMessage = "抢购成功"
		}
	}
	this.Ctx.WriteString(showMessage)
}
