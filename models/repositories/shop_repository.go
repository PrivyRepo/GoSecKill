package repositories

import (
	"database/sql"
	"errors"
	"homework/common/mysql"
	"homework/common/reflect"
	"homework/models/datamodels"
	"strconv"
)

type IShopRepository interface {
	Conn() error
	Select(UserNmae string) (*datamodels.Shop, error)
	Insert(user *datamodels.Shop) (shopId int64, err error)
	CloseConn()
}

type ShopRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (s *ShopRepository) CloseConn() {
	s.mysqlConn.Close()
}

func (s *ShopRepository) Conn() error {
	if s.mysqlConn == nil {
		db, err := mysql.NewMysqlConn()
		if err != nil {
			return err
		}
		s.mysqlConn = db
	}
	if s.table == "" {
		s.table = "shop"
	}
	return nil
}

func (s *ShopRepository) Select(UserName string) (*datamodels.Shop, error) {
	if UserName == "" {
		return &datamodels.Shop{}, errors.New("不能为空")
	}
	if err := s.Conn(); err != nil {
		return &datamodels.Shop{}, err
	}
	sql := "SELECT * FROM " + s.table + " where userName=?"
	row, e := s.mysqlConn.Query(sql, UserName)
	defer row.Close()
	if e != nil {
		return nil, e
	}
	shop := &datamodels.Shop{}
	resultRow := mysql.GetResultRow(row)
	if len(resultRow) == 0 {
		return &datamodels.Shop{}, errors.New("店铺不存在")
	}
	reflect.DataToStructByTagSql(resultRow, shop)
	return shop, nil
}

func (s *ShopRepository) Insert(shop *datamodels.Shop) (userId int64, err error) {
	if err = s.Conn(); err != nil {
		return
	}
	sql := "INSERT " + s.table + " SET shopName=?,userName=?,passWord=?"
	stmt, err := s.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return
	}
	result, err := stmt.Exec(shop.ShopName, shop.UserName, shop.HashPassword)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (s *ShopRepository) SelectByID(userId int64) (*datamodels.Shop, error) {
	sql := "SELECT * FROM " + s.table + "WHERE ID =" + strconv.FormatInt(userId, 10)
	rows, e := s.mysqlConn.Query(sql)
	defer rows.Close()
	if e != nil {
		return nil, e
	}
	row := mysql.GetResultRow(rows)
	shop := &datamodels.Shop{}
	reflect.DataToStructByTagSql(row, shop)
	return shop, nil
}

func NewShopRepository(db *sql.DB, table string) IShopRepository {
	return &ShopRepository{table, db}
}
