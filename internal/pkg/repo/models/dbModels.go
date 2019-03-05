package models

import (
	"github.com/shopspring/decimal"
)

type Holder struct {
	ID         uint `gorm:"primary_key"`
	EthAddress string
	Balance    decimal.Decimal `gorm:"type:decimal"`
	Ok         bool
}

type HolderNew struct {
	ID         uint `gorm:"primary_key"`
	EthAddress string
	Balance    decimal.Decimal `gorm:"type:decimal"`
	Ok         bool
}
