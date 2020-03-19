package repositories

//第一步,先开发对应的接口
//第二步,实现接口
import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"homework/common"
	"homework/models/datamodels"
	"strconv"
)

//第一步，先开发对应的接口
//第二步，实现定义的接口
type IProduct interface {
	//连接数据
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]datamodels.Product, error)
	SelectByshopId(int64) ([]datamodels.Product, error)
	SubProductNum(productID int64) error
	IncrProductReview(productId int64) bool
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
	redisConn redis.Conn
}

func NewProductManager(table string, db *sql.DB, redis redis.Conn) IProduct {
	return &ProductManager{table: table, mysqlConn: db, redisConn: redis}
}

//数据连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

func (p *ProductManager) RedisConn() (err error) {
	if p.redisConn == nil {
		conn := common.NewRedisConn()
		p.redisConn = conn
	}
	return
}

//插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	//1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	//2.准备sql
	sql := "INSERT product SET shopID=?,productName=?,productNum=?,productImage=?,productReviews=?,productOldprice=?,productNewprice=?,productDescription=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}
	//3.传入参数
	result, errStmt := stmt.Exec(product.Shopid, product.ProductName, product.ProductNum, product.ProductImage, product.ProductReviews, product.ProductOldprice, product.ProductNewprice, product.ProductDescription)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

//商品的删除
func (p *ProductManager) Delete(productID int64) bool {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
	if err != nil {
		return false
	}
	return true
}

//商品的更新
func (p *ProductManager) Update(product *datamodels.Product) error {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "Update product set productName=?,productNum=?,productImage=?,productReviews=?,productOldprice=?,productNewprice=?,productDescription=? where ID = " + strconv.FormatInt(product.ID, 10)

	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductReviews, product.ProductOldprice, product.ProductNewprice, product.ProductDescription)
	if err != nil {
		return err
	}
	return nil
}

//根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	//1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return

}

//获取所有商品
func (p *ProductManager) SelectAll() (productArray []datamodels.Product, errProduct error) {

	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil, err
	}

	rel, _ := redis.Bytes(p.redisConn.Do("get", "products"))

	decoder := gob.NewDecoder(bytes.NewReader(rel))
	decoder.Decode(&productArray)

	if len(productArray) == 0 {
		fmt.Println("数据不存在redis，从数据库中获取")
		//1.判断连接是否存在
		if err := p.Conn(); err != nil {
			return nil, err
		}
		sql := "Select * from " + p.table + " ORDER BY ID "
		rows, errProduct := p.mysqlConn.Query(sql)
		defer rows.Close()
		if errProduct != nil {
			return nil, errProduct
		}

		result := common.GetResultRows(rows)
		if len(result) == 0 {
			return nil, nil
		}

		for _, v := range result {
			product := &datamodels.Product{}
			common.DataToStructByTagSql(v, product)
			productArray = append(productArray, *product)
		}
		//var buffer bytes.Buffer
		//encoder := gob.NewEncoder(&buffer)
		//_ = encoder.Encode(productArray)
		//_, err := p.redisConn.Do("set", "products", buffer.Bytes(),"EX","1")
		//if err != nil {
		//	fmt.Println(err)
		//}
	} else {
		fmt.Println("数据存在redis，从redis中获取")
	}
	return
}
func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update " + p.table + " set " + " productNum=productNum-1 where ID =" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func (p *ProductManager) SelectByshopId(shopid int64) (productArray []datamodels.Product, errProduct error) {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil, err
	}

	sql := "Select * from " + p.table + " WHERE shopID = " + strconv.FormatInt(shopid, 10) + " ORDER BY ID "
	rows, errProduct := p.mysqlConn.Query(sql)
	defer rows.Close()
	if errProduct != nil {
		return nil, errProduct
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, *product)
	}
	return
}

func (p *ProductManager) IncrProductReview(productId int64) bool {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "UPDATE " + p.table + " SET productreviews = productreviews + 1  WHERE ID = " + strconv.FormatInt(productId, 10)
	_, e := p.mysqlConn.Exec(sql)
	if e != nil {
		return false
	}
	return true
}
