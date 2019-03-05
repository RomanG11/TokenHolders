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
	FindGroup(st, f int64) ([]models.Holder, error)
	FindAll() ([]models.Holder, error)
	FindAllWithPositiveBalance() ([]models.Holder, error)
	FindHolder(addr string) (models.Holder, error)
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
		Balance:    balance,
	}

	err := repo.db.Create(&holder).Error
	if err != nil {
		return holder, err
	}

	return holder, nil
}

func (repo *HolderRepo) FindGroup(st, f int64) ([]models.Holder, error) {

	var h []models.Holder

	err := repo.db.Where("id > ? and id < ? ", st, f).Find(&h).Error
	if err != nil {
		return h, err
	}

	return h, nil
}

func (repo *HolderRepo) FindAll() ([]models.Holder, error) {

	var h []models.Holder

	err := repo.db.Find(&h).Error
	if err != nil {
		return h, err
	}

	return h, nil
}

func (repo *HolderRepo) FindAllWithPositiveBalance() ([]models.Holder, error) {

	var h []models.Holder

	err := repo.db.Where("balance>0 and ok is null").Find(&h).Error

	if err != nil {
		return h, err
	}

	return h, nil
}

func (repo *HolderRepo) FindHolder(addr string) (models.Holder, error) {

	var h models.Holder

	err := repo.db.Where("eth_address = ?", addr).Find(&h).Error
	if err != nil {
		return h, err
	}

	return h, nil
}
