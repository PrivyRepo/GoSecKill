package services

import (
	"homework/models/datamodels"
	"homework/models/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	GetOrderInfoByID(int64) (map[string]string, error)
	UpdateOrderByID(int64) bool
	InsertOrderByMessage(*datamodels.Message) (int64, error)
}

type OrderService struct {
	orderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID(id int64) (*datamodels.Order, error) {
	return o.orderRepository.SelectByKey(id)
}

func (o *OrderService) DeleteOrderByID(id int64) bool {
	return o.orderRepository.Delete(id)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.orderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.orderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.orderRepository.SelectAllWithInfo()
}

func (o *OrderService) GetOrderInfoByID(id int64) (map[string]string, error) {
	return o.orderRepository.SelectWithInfoByKey(id)
}

func (o *OrderService) UpdateOrderByID(id int64) bool {
	return o.orderRepository.UpdateInfoByKey(id)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{repository}
}

//根据消息创建订单
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
