package rabbitmq

import (
	"database/sql"
	"fmt"
	"github.com/streadway/amqp"
	"homework/models/datamodels"
	"log"

	"encoding/json"
	"sync"
)

//连接信息
//const MQURL = "amqp://imoocuser:imoocuser@172.31.96.59:5672/imooc"
const MQURL = "amqp://imoocuser:imoocuser@192.168.19.200:5672/imooc"

//简单的消息体
type Message struct {
	ProductID int64
	UserID    int64
}

//创建结构体
func NewMessage(userId int64, productId int64) *Message {
	return &Message{productId, userId}
}

//rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
	sync.Mutex
}

//创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		return err
	}
	//调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple(db *sql.DB) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//消费者流控
	r.channel.Qos(
		1,     //当前消费者一次能接受的最大消息数量
		0,     //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		//这里要改掉，我们用手动应答
		false, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			message := &Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			tx, e := db.Begin()
			if e != nil {
				updateSQL := "UPDATE `product` SET productNum = productNum -1 WHERE ID = ?"
				_, e2 := tx.Exec(updateSQL, message.ProductID)
				insertSQL := "INSERT `order` set userID=?,productID=?,orderPayStatus=?,orderDeliverStatus=?"
				_, e1 := tx.Exec(insertSQL, message.UserID, message.ProductID, datamodels.PayWait, datamodels.DeliverWait)
				if e1 != nil || e2 != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}
			//如果为true表示确认所有未确认的消息，
			//为false表示确认当前消息
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
