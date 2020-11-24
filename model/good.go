package model

type Good struct {
	GoodId		uint `json:"good_id" gorm:"primary_key"`
	UserId		uint `json:"user_id" gorm:"type:int;not null"`
	Quantity	uint `json:"quantity" gorm:"type:int;not null"` // 商品库存数量
	Name	 	string `json:"name" gorm:"type:varchar(20);not null"`
	Price		float32 `json:"price" gorm:"type:float"`
	Info		string `json:"info" gorm:"type:text"`
	Picture		string `json:"img" gorm:"type:varchar(128)"`
}