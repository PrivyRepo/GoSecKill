package datamodels

type Shop struct {
	ID           int64  `json:"id" form:"ID" sql:"ID"`
	ShopName     string `json:"shopName" form:"shopName" sql:"shopName"`
	UserName     string `json:"userName" form:"userName" sql:"userName"`
	HashPassword string `json:"-" form:"passWord" sql:"passWord"`
}
