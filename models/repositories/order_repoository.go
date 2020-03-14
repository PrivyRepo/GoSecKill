package repositories

import (
	"database/sql"
	"homework/common"
	"homework/models/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(product *datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
	SelectWithInfoByKey(int64) (map[string]string, error)
	UpdateInfoByKey(int64) bool
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table, sql}

}

type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) UpdateInfoByKey(id int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "UPDATE `order` SET orderStatus = 1 WHERE ID = " + strconv.FormatInt(id, 10)
	_, e := o.mysqlConn.Exec(sql)
	if e != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT `" + o.table + "` set userID=?,productID=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return
	}
	productID, err = result.LastInsertId()
	return

}

func (o *OrderManagerRepository) Delete(productid int64) (isok bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sql := "DELETE FROM `" + o.table + "` WHERE ID=?"
	stmt, e := o.mysqlConn.Prepare(sql)
	if e != nil {
		return false
	}
	_, e = stmt.Exec(productid)
	if e != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "Update " + o.table + "SET userID=?,productID=?,orderStatus=? WHERE ID = " + strconv.FormatInt(order.ID, 10)
	stmt, e := o.mysqlConn.Prepare(sql)
	if e != nil {
		return e
	}
	_, e = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.ID)
	return e

}

func (o *OrderManagerRepository) SelectByKey(id int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "SELECT * FROM" + o.table + "WHERE ID = " + strconv.FormatInt(id, 10)
	rows, e := o.mysqlConn.Query(sql)
	if e != nil {
		return &datamodels.Order{}, e
	}
	resultRows := common.GetResultRow(rows)
	if len(resultRows) == 0 {
		return &datamodels.Order{}, e
	}
	order := &datamodels.Order{}
	common.DataToStructByTagSql(resultRows, order)
	return order, nil

}

func (o *OrderManagerRepository) SelectAll() ([]*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT * FROM" + o.table
	rows, e := o.mysqlConn.Query(sql)
	if e != nil {
		return nil, e
	}
	resultRows := common.GetResultRows(rows)
	if len(resultRows) == 0 {
		return nil, nil
	}
	res := make([]*datamodels.Order, 0)
	for _, v := range resultRows {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		res = append(res, order)
	}
	return res, nil

}

func (o *OrderManagerRepository) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT o.ID,p.productName,u.userName,o.orderStatus FROM `order` as o LEFT JOIN product as p on o.productID=p.ID JOIN user as u ON o.userID = u.ID"
	rows, e := o.mysqlConn.Query(sql)
	if e != nil {
		return nil, e
	}
	return common.GetResultRows(rows), nil
}

func (o *OrderManagerRepository) SelectWithInfoByKey(id int64) (map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT o.ID , u.userName, p.productName, o.orderStatus FROM `order` as o JOIN product as p ON o.productID = p.ID JOIN user as u ON o.userID = u.ID"
	rows, e := o.mysqlConn.Query(sql)
	if e != nil {
		return nil, e
	}
	return common.GetResultRow(rows), nil
}
