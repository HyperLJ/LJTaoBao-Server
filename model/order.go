package model

import "github.com/jinzhu/gorm"

// 购物车
type Order struct {
	gorm.Model
	UserId 		uint 	`json:"user_id" gorm:"type:int;not null"`
	GoodId		uint 	`json:"good_id" gorm:"type:int;not null"`
	GoodsCount	uint 	`json:"goods_count" gorm:"type:int;not null"`
	GoodsPrice	float32 `json:"goods_price" gorm:"type:float;not null"`//商品总价
}