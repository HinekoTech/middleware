package orm

import "github.com/FourWD/middleware/model"

type UserProfile struct {
	ID string `json:"id" query:"id" gorm:"type:varchar(36);primary_key"`
	model.GormModel

	UserID                string `json:"user_id" query:"user_id" gorm:"type:varchar(36)"`
	Address       string `json:"address" query:"address" gorm:"type:text"`
	Street        string `json:"street" query:"street" gorm:"type:varchar(200)"`
	DistrictID    string `json:"district_id" query:"district_id" gorm:"type:varchar(4)"`         //อำเภอ
	SubDistrictID string `json:"sub_district_id" query:"sub_district_id" gorm:"type:varchar(6)"` //ตำบล
	ProvinceID    string `json:"province_id" query:"province_id" gorm:"type:varchar(2)"`
}
