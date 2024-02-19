package orm

import (
	"time"

	"github.com/FourWD/middleware/model"
)

type CallbackTime struct {
	ID string `json:"id" query:"id" gorm:"type:varchar(36);primary_key"`
	model.GormModel

	StartDate time.Time `json:"start_date" query:"start_date"`
	EndDate   time.Time `json:"end_date" query:"end_date"`
	Name      string    `json:"name" query:"name" gorm:"not null;type:varchar(50)"`

	model.GormRowOrder
}
