package models

import "gorm.io/gorm"

type Rate struct {
	gorm.Model
	BaseCode string  `gorm:"uniqueIndex:idx_basecode_code"`
	Code     string  `gorm:"uniqueIndex:idx_basecode_code"`
	Rate     float64 `json:"rate"`
}

func (r *Rate) TableName() string {
	return "rates"
}
