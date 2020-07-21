package repositories

import (
	"database/sql"
	"fmt"
	mysql2 "homework/common/mysql"
	"homework/common/reflect"
	"homework/models/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(product *datamodels.Order) (int, error)
	Delete(int) bool
	Update(order *datamodels.Order) error
	SelectById(int) (order *datamodels.Order, err error)
	SelectWithInfoByUser(int, int, int) (map[int]map[string]string, int, error)
	SelectWithInfoByShop(int, int, int) (map[int]map[string]string, int, error)
	CloseConn()
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table, sql}
}

type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) CloseConn() {
	o.mysqlConn.Close()
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := mysql2.NewMysqlConn()
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


func (o *OrderManagerRepository) SelectById(id int) (order *datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "SELECT * FROM `order` WHERE id = " + strconv.Itoa(id)
	row, e := o.mysqlConn.Query(sql)
	defer row.Close()
	if e != nil {
		return nil, e
	}
	result := mysql2.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}
	orderResult := &datamodels.Order{}
	reflect.DataToStructByTagSql(result, orderResult)
	return orderResult, nil

}
func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT `order` set userID=?,productID=?,orderPayStatus=?,orderDeliverStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderPayStatus, order.OrderDeliverStatus)
	if err != nil {
		return
	}
	productid, err := result.LastInsertId()
	return int(productid), err

}

func (o *OrderManagerRepository) Delete(productid int) (isok bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sql := "DELETE FROM `order` WHERE ID=?"
	stmt, e := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
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
	sql := "Update `order` SET userID=?,productID=?,orderPayStatus=?,orderDeliverStatus=? WHERE ID = " + strconv.FormatInt(order.ID, 10)
	stmt, e := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if e != nil {
		return e
	}
	_, e = stmt.Exec(order.UserId, order.ProductId, order.OrderPayStatus, order.OrderDeliverStatus)
	return e
}

func (o *OrderManagerRepository) SelectWithInfoByUser(id int, start int, limit int) (map[int]map[string]string, int, error) {
	if err := o.Conn(); err != nil {
		return nil, 0, err
	}
	var count int
	CountSql := "SELECT COUNT(*) FROM `order` WHERE userID = " + strconv.Itoa(id)
	row := o.mysqlConn.QueryRow(CountSql)
	err := row.Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	sql := "SELECT o.ID , s.shopName, p.productName, o.orderPayStatus,o.orderDeliverStatus " +
		"FROM `order` as o " +
		"JOIN product as p ON o.productID = p.ID " +
		"JOIN user as u ON o.userID = u.ID " +
		"JOIN shop as s ON p.shopID = s.ID " +
		"WHERE u.id = " + strconv.Itoa(id) + fmt.Sprintf(" LIMIT %d,%d", start, limit)
	rows, e := o.mysqlConn.Query(sql)
	defer rows.Close()
	if e != nil {
		return nil, count, e
	}
	return mysql2.GetResultRows(rows), count, nil
}

func (o *OrderManagerRepository) SelectWithInfoByShop(id int, start int, limit int) (map[int]map[string]string, int, error) {
	if err := o.Conn(); err != nil {
		return nil, 0, err
	}
	var count int
	CountSql := "SELECT COUNT(*) FROM `order` as o " +
		"JOIN `product` as p ON o.productID = p.ID " +
		"JOIN `shop` as s ON p.shopID = s.ID" +
		" WHERE s.ID = " + strconv.Itoa(id)
	row := o.mysqlConn.QueryRow(CountSql)
	err := row.Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	sql := "SELECT o.ID , u.userName, p.productName, o.orderPayStatus,o.orderDeliverStatus " +
		"FROM `order` as o " +
		"JOIN product as p ON o.productID = p.ID " +
		"JOIN user as u ON o.userID = u.ID " +
		"JOIN shop as s ON p.shopID = s.ID " +
		"WHERE s.id = " + strconv.Itoa(id) + fmt.Sprintf(" LIMIT %d,%d", start, limit)
	rows, e := o.mysqlConn.Query(sql)
	defer rows.Close()
	if e != nil {
		return nil, count, e
	}
	return mysql2.GetResultRows(rows), count, nil
}
