package services

import (
	"homework/models/datamodels"
	"homework/models/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(productID int64) error
	IncPorductReview(productId int64) bool
	GetProductByshop(shopID int64) ([]datamodels.Product, error)
}

type ProductService struct {
	productRepository repositories.IProduct
}

//初始化函数
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *ProductService) GetProductByshop(shopID int64) ([]datamodels.Product, error) {
	return p.productRepository.SelectByshopId(shopID)
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}

func (p *ProductService) IncPorductReview(productId int64) bool {
	return p.productRepository.IncrProductReview(productId)
}
