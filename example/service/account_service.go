package service

import (
	"github.com/henrjan/lock_with_key/example/model"
	"github.com/henrjan/lock_with_key/example/repository"
)

type accountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *accountService {
	return &accountService{repo: repo}
}

func (service *accountService) GetAllAccount() []model.Account {
	return service.repo.GetAllAccount()
}

func (service *accountService) GetAccountByNumber(accNumber string) model.Account {
	return service.repo.GetAccountByNumber(accNumber)
}
