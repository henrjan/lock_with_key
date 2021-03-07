package repository

import (
	"github.com/henrjan/lock_with_key/example/model"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(dbClient *gorm.DB) *AccountRepository {
	dbClient.AutoMigrate(&model.Account{})
	dbClient.Where("1 = 1").Delete(&model.Account{})

	dbClient.Create(&model.AccountList)
	return &AccountRepository{dbClient}
}

func (repo *AccountRepository) InsertMany(model []model.Account) {
	repo.db.Create(&model)
}

func (repo *AccountRepository) GetAllAccount() []model.Account {
	accountList := make([]model.Account, 0)

	repo.db.Find(&accountList)
	return accountList
}

func (repo *AccountRepository) GetAccountByNumber(accNumber string) model.Account {
	accountList := model.Account{}

	repo.db.Limit(1).Where("account_number = ?", accNumber).Find(&accountList)
	return accountList
}
