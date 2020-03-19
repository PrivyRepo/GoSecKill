package services

import (
	"homework/models/datamodels"
	"homework/models/repositories"
)

type IShopService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.Shop, isOk bool)
	AddShop(shop *datamodels.Shop) (shopId int64, err error)
}

type ShopService struct {
	ShopRepository repositories.IShopRepository
}

func NewShopService(repository repositories.IShopRepository) IShopService {
	return &ShopService{repository}
}

func (s *ShopService) IsPwdSuccess(userName string, pwd string) (shop *datamodels.Shop, isOK bool) {
	shop, e := s.ShopRepository.Select(userName)
	if e != nil {
		return
	}
	isOK, _ = ValidatePassword(pwd, shop.HashPassword)
	if !isOK {
		return &datamodels.Shop{}, false
	}
	return
}

func (s *ShopService) AddShop(shop *datamodels.Shop) (shopId int64, err error) {
	bytes, err := GeneratePassword(shop.HashPassword)
	if err != nil {
		return shopId, err
	}
	shop.HashPassword = string(bytes)
	return s.ShopRepository.Insert(shop)
}
