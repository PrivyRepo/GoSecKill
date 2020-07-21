package main

import (
	"fmt"
	"homework/common/mysql"
	"homework/seckill/tool/config"
	"homework/seckill/tool/consistent"
	"io/ioutil"
	"net/http"
)

//定时
//清零，预热

var hashConsistent *consistent.Consistent

func init(){
	hostArray := config.ReadLineFile("appsconfig")
	hashConsistent = consistent.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
}


func main(){
	db, err := mysql.NewMysqlConn()
	if err!=nil {
		panic("mysql connect error")
	}
	sql := "SELECT ID , productNum FROM product WHERE productNum > 0"
	rows, err := db.Query(sql)
	if err!=nil {
		panic("mysql query error")
	}
	resultRows := mysql.GetResultRows(rows)
	client := &http.Client{}
	for _,v := range resultRows {
		host, _ := hashConsistent.Get(v["ID"])
		resp, err := client.Get(fmt.Sprintf("http://%s:8083/seckill/setCount?productid=%s&count=%s", host, v["ID"], v["productNum"]))
		if err!=nil {
			fmt.Println(err)
		}
		bs,_ := ioutil.ReadAll(resp.Body)
		if string(bs) != "true" {
			fmt.Println(resp.Header)
		}
	}
	fmt.Println("配置完成")
}
