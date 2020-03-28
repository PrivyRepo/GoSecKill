package services

import (
	"homework/common/rabbitmq"
	"homework/models/datamodels"
	"homework/models/repositories"
)

type IOrderService interface {
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int, error)
	GetOrderByID(int) (*datamodels.Order, error)
	GetOrderInfoByUser(int, int, int) (map[int]map[string]string, int, error)
	GetOrderInfoByShop(int, int, int) (map[int]map[string]string, int, error)
	InsertOrderByMessage(*rabbitmq.Message) (int, error)
}

type OrderService struct {
	orderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID(id int) (*datamodels.Order, error) {
	return o.orderRepository.SelectById(id)
}
func (o *OrderService) DeleteOrderByID(id int) bool {
	return o.orderRepository.Delete(id)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int, error) {
	return o.orderRepository.Insert(order)
}

func (o *OrderService) GetOrderInfoByUser(id int, pagenum int, limit int) (map[int]map[string]string, int, error) {
	return o.orderRepository.SelectWithInfoByUser(id, (pagenum-1)*10, limit)
}

func (o *OrderService) GetOrderInfoByShop(id int, pagenum int, limit int) (map[int]map[string]string, int, error) {
	return o.orderRepository.SelectWithInfoByShop(id, (pagenum-1)*10, limit)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{repository}
}

//根据消息创建订单
func (o *OrderService) InsertOrderByMessage(message *rabbitmq.Message) (orderID int, err error) {
	order := &datamodels.Order{
		UserId:         message.UserID,
		ProductId:      message.ProductID,
		OrderPayStatus: datamodels.PayWait,
	}
	return o.InsertOrder(order)
}
