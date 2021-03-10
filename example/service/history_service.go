package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/henrjan/lock_with_key/example/general"
	"github.com/henrjan/lock_with_key/example/model"
	"github.com/henrjan/lock_with_key/example/repository"
)

var (
	transferFee = 500
	minTransfer = 10000
)

type historyService struct {
	repo            *repository.HistoryRepository
	mLock           *general.MultiLock
	dataTransaction map[string][]model.HistoryLog
}

func NewHistoryService(repo *repository.HistoryRepository, mLock *general.MultiLock) *historyService {
	return &historyService{
		repo:            repo,
		mLock:           mLock,
		dataTransaction: make(map[string][]model.HistoryLog),
	}
}

func (service *historyService) GetAccountLastBalance(from model.Account) uint {
	service.mLock.RLock(from.AccountNumber)
	defer service.mLock.RUnlock(from.AccountNumber)
	txHistory := service.repo.GetUserLatestTransaction(from.ID)
	return txHistory.CurrentBalance
}

func (service *historyService) GetAccountTxHistories(from model.Account) []model.HistoryLog {
	service.mLock.RLock(from.AccountNumber)
	defer service.mLock.RUnlock(from.AccountNumber)
	return service.repo.GetAccountTxHistories(from.ID)
}

func (service *historyService) TransferBalance(from model.Account, to model.Account, nominal uint) {
	referenceID := uuid.New()

	service.SubtractBalance(from, to, nominal, referenceID)
	service.AddBalance(from, to, nominal, referenceID)
	service.SubtractBalance(from, model.Account{}, uint(transferFee), referenceID)
}

func (service *historyService) AddBalance(from model.Account, to model.Account, nominal uint, transactionID uuid.UUID) {
	service.mLock.Lock(to.AccountNumber)
	defer service.mLock.Unlock(to.AccountNumber)
	getLatestTransaction := service.repo.GetUserLatestTransaction(to.ID)
	balanceAfterTx := getLatestTransaction.CurrentBalance + nominal

	newTransaction := model.HistoryLog{
		SenderID:       from.ID,
		OwnerID:        to.ID,
		LastBalance:    getLatestTransaction.CurrentBalance,
		Balance:        int(nominal),
		CurrentBalance: balanceAfterTx,
		Tag:            "BALANCE_IN",
		CreatedAt:      time.Now(),
		ReferenceID:    transactionID,
	}
	service.repo.Insert(newTransaction)
}

func (service *historyService) SubtractBalance(from model.Account, to model.Account, nominal uint, transactionID uuid.UUID) {
	service.mLock.Lock(from.AccountNumber)
	defer service.mLock.Unlock(from.AccountNumber)
	var latestBalance uint

	getLatestTransaction := service.repo.GetUserLatestTransaction(from.ID)
	latestBalance = getLatestTransaction.CurrentBalance
	balanceAfterTx := latestBalance - nominal

	newTransaction := model.HistoryLog{
		OwnerID:        from.ID,
		RecipientID:    to.ID,
		LastBalance:    latestBalance,
		Balance:        int(nominal) * -1,
		CurrentBalance: balanceAfterTx,
		Tag:            "BALANCE_OUT",
		CreatedAt:      time.Now(),
		ReferenceID:    transactionID,
	}
	service.repo.Insert(newTransaction)
}
