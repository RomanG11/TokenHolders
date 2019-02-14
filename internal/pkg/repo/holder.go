package repo

import (
	"TokenHolders/internal/pkg/repo/models"
	"github.com/shopspring/decimal"

	"github.com/jinzhu/gorm"
)

type HolderRepository interface {
	GetHolderByAddress(ethAddress string) (models.Holder, error)
	UpdateHolder(holder *models.Holder) error
	NewHolder(address string, balance decimal.Decimal) (models.Holder, error)
}

type HolderRepo struct {
	db *gorm.DB
}

func (repo *HolderRepo) GetHolderByAddress(ethAddress string) (models.Holder, error) {
	var holder models.Holder

	err := repo.db.Where("eth_address = ?", ethAddress).Find(&holder).Error
	if err != nil {
		return holder, err
	}

	return holder, nil
}

func (repo *HolderRepo) UpdateHolder(holder *models.Holder) error {

	err := repo.db.Save(holder).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *HolderRepo) NewHolder(address string, balance decimal.Decimal) (models.Holder, error) {

	holder := models.Holder{
		EthAddress: address,
		Balance: balance,
	}

	err := repo.db.Create(&holder).Error
	if err != nil {
		return holder, err
	}

	return holder, nil
}

func (repo *HolderRepo) FindGroup(st, f int64) ([]models.Holder, error) {
	repo.db.Where("id  ? and ", )
}