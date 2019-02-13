package models

import (
	"github.com/shopspring/decimal"
)

type Holder struct {
	EthAddress          string               `gorm:"primary_key"`
	Balance			    decimal.Decimal      `gorm:"type:decimal"`
}
