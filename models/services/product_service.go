package services

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"homework/models/datamodels"
	"homework/models/repositories"
	"log"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct(int, int) ([]datamodels.Product, int, error)
	GetAllProductInfo(int, int) (map[int]map[string]string, int, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(productID int64) error
	GetProductByshop(shopID int64, pagenum int, limit int) ([]datamodels.Product, int, error)
}

type ProductService struct {
	productRepository repositories.IProduct
}

//初始化函数
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

func (p *ProductService) GetAllProductInfo(pagenum int, limit int) (map[int]map[string]string, int, error) {
	return p.productRepository.SelectAllInfo((pagenum-1)*limit, limit)
}

func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct(pagenum int, limit int) ([]datamodels.Product, int, error) {
	return p.productRepository.SelectAll((pagenum-1)*limit, limit)
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	if err := ValidateProduct(product); err != nil {
		return -1, err
	}
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *ProductService) GetProductByshop(shopID int64, pagenum int, limit int) ([]datamodels.Product, int, error) {
	return p.productRepository.SelectByshopId(shopID, (pagenum-1)*10, limit)
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}

func ValidateProduct(product *datamodels.Product) error {
	valid := validation.Validation{}
	b, err := valid.Valid(product)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		err = errors.New("数据验证错误")
	}
	return err
}
