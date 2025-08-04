package models

import (
	"gorm.io/gorm"
	"time"
)

type PaymentMethod struct {
	ID                uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UID               string         `gorm:"not null" json:"uid"`
	UserId            *string        `json:"userId"`                            // NULL for default rates, value for custom rates
	PaymentMethodType uint           `gorm:"not null" json:"paymentMethodType"` // references consts.PaymentMethodsType
	CurrencyFrom      uint           `gorm:"not null" json:"currencyFrom"`
	CurrencyTo        uint           `gorm:"not null" json:"currencyTo"`
	Rate              float64        `json:"rate"`   // TASA
	Margin            float64        `json:"margin"` // MARGEN
	Fee               float64        `json:"fee"`    // COMISION
	Name              string         `gorm:"not null" json:"name"`
	Description       string         `gorm:"not null" json:"description"`
	IsActive          bool           `gorm:"not null;default:true" json:"isActive"`
	IsCustom          bool           `json:"isCustom"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

type PaymentMethodsRequest struct {
	CurrencyFrom uint   `json:"currencyFrom"`
	CurrencyTo   uint   `json:"currencyTo"`
	UserId       string `json:"userId"`
}

// TableName specifies the table name
func (pm *PaymentMethod) TableName() string {
	return "payment_methods"
}
