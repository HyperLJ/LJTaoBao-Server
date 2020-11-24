package model

//// gorm.Model 定义
//type Model struct {
//	ID        uint `gorm:"primary_key"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt *time.Time
//}

type User struct {
	UserId		uint	`json:"user_id" gorm:"primary_key"`
	Account		string `json:"account" gorm:"type:varchar(20);not null"`
	Password	string `json:"password" gorm:"type:varchar(20);not null"`
	Name	 	string `json:"name" gorm:"type:varchar(20);not null"`
	Sex			bool `json:"sex" gorm:"type:bool;default:true"`
	Info 		string `json:"info" gorm:"type:text"`
	Head		string `json:"head" gorm:"type:varchar(128)"`
}