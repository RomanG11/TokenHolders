package repo

import (
	"TokenHolders/internal/pkg/repo/models"
	"github.com/shopspring/decimal"

	"github.com/jinzhu/gorm"
)

type HolderNewRepository interface {
	GetHolderByAddress(ethAddress string) (models.HolderNew, error)
	UpdateHolder(holder *models.HolderNew) error
	NewHolder(address string, balance decimal.Decimal) (models.HolderNew, error)
	FindGroup(st, f int64) ([]models.HolderNew, error)
	FindAll() ([]models.HolderNew, error)
	FindAllWithPositiveBalance() ([]models.HolderNew, error)
	FindHolder(addr string) (models.HolderNew, error)
}

type HolderNewRepo struct {
	db *gorm.DB
}

func (repo *HolderNewRepo) GetHolderByAddress(ethAddress string) (models.HolderNew, error) {
	var holder models.HolderNew

	err := repo.db.Where("eth_address = ?", ethAddress).Find(&holder).Error
	if err != nil {
		return holder, err
	}

	return holder, nil
}

func (repo *HolderNewRepo) UpdateHolder(holder *models.HolderNew) error {

	err := repo.db.Save(holder).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *HolderNewRepo) NewHolder(address string, balance decimal.Decimal) (models.HolderNew, error) {

	holder := models.HolderNew{
		EthAddress: address,
		Balance:    balance,
	}

	err := repo.db.Create(&holder).Error
	if err != nil {
		return holder, err
	}

	return holder, nil
}

func (repo *HolderNewRepo) FindGroup(st, f int64) ([]models.HolderNew, error) {

	var h []models.HolderNew

	err := repo.db.Where("id > ? and id < ? ", st, f).Find(&h).Error
	if err != nil {
		return h, err
	}

	return h, nil
}

func (repo *HolderNewRepo) FindAll() ([]models.HolderNew, error) {

	var h []models.HolderNew

	err := repo.db.Find(&h).Error
	if err != nil {
		return h, err
	}

	return h, nil
}

func (repo *HolderNewRepo) FindAllWithPositiveBalance() ([]models.HolderNew, error) {

	var h []models.HolderNew

	err := repo.db.Where("balance>0 and ok is null").Find(&h).Error

	if err != nil {
		return h, err
	}

	return h, nil
}

func (repo *HolderNewRepo) FindHolder(addr string) (models.HolderNew, error) {

	var h models.HolderNew

	err := repo.db.Where("eth_address = ?", addr).Find(&h).Error
	if err != nil {
		return h, err
	}

	return h, nil
}
