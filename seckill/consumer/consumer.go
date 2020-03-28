package main

import (
	"fmt"
	"homework/common/mysql"
	"homework/common/rabbitmq"
)

func main() {
	db, err := mysql.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("imoocProduct")
	rabbitmqConsumeSimple.ConsumeSimple(db)
}
