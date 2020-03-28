package datamodels

type Product struct {
	ID                 int64  `json:"id" sql:"ID" form:"id"`
	Shopid             int64  `json:"shopid" sql:"shopID" form:"shopid" "`
	ProductName        string `json:"ProductName" sql:"productName" form:"ProductName" `
	ProductNum         int64  `json:"ProductName" sql:"productNum" form:"ProductNum" "`
	ProductImage       string `json:"ProductImage" sql:"productImage" form:"ProductImage" `
	ProductOldprice    int    `json:"ProductOldprice" sql:"productOldprice" form:"ProductOldprice" valid:"Range(1, 100000000000000000)`
	ProductNewprice    int    `json:"ProductNewprice" sql:"productNewprice" form:"ProductNewprice" valid:"Range(1, 100000000000000000)`
	ProductDescription string `json:"ProductDescription" sql:"productDescription" form:"ProductDescription" `
}
