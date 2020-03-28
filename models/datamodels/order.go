package datamodels

type Order struct {
	ID                 int64 `sql:"ID"`
	UserId             int64 `sql:"userID"`
	ProductId          int64 `sql:"productID"`
	OrderPayStatus     int   `sql:"orderPayStatus"`
	OrderDeliverStatus int   `sql:"orderDeliverStatus"`
}

const (
	PayWait    = iota
	PaySuccess //1
	PayFailed  //2
)

const (
	DeliverWait    = iota
	DeliverSuccess //1
)
