package repositories

import (
	"database/sql"
	"errors"
	"github.com/gomodule/redigo/redis"
	"homework/common"
	"homework/models/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(UserName string) (*datamodels.User, error)
	Insert(user *datamodels.User) (userId int64, err error)
}

type UserRepository struct {
	table     string
	mysqlConn *sql.DB
	redisConn redis.Conn
}

func (u UserRepository) Conn() error {
	if u.mysqlConn == nil {
		db, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = db
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u UserRepository) Select(UserName string) (*datamodels.User, error) {
	if UserName == "" {
		return &datamodels.User{}, errors.New("不能为空")
	}
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "SELECT * FROM " + u.table + " where UserName=?"
	row, e := u.mysqlConn.Query(sql, UserName)
	if e != nil {
		return nil, e
	}
	user := &datamodels.User{}
	resultRow := common.GetResultRow(row)
	if len(resultRow) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	common.DataToStructByTagSql(resultRow, user)
	return user, nil
}

func (u UserRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (u *UserRepository) SelectByID(userId int64) (*datamodels.User, error) {
	sql := "SELECT * FROM " + u.table + "WHERE ID =" + strconv.FormatInt(userId, 10)
	rows, e := u.mysqlConn.Query(sql)
	if e != nil {
		return nil, e
	}
	row := common.GetResultRow(rows)
	user := &datamodels.User{}
	common.DataToStructByTagSql(row, user)
	return user, nil
}

func NewUserRepository(table string, db *sql.DB, redisConn redis.Conn) IUserRepository {
	return &UserRepository{table, db, redisConn}
}
