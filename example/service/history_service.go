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
	dataTransaction []*model.HistoryLog
	cacheLatestTx   map[string]*model.HistoryLog
}

func NewHistoryService(repo *repository.HistoryRepository, mLock *general.MultiLock) *historyService {
	return &historyService{
		repo:          repo,
		mLock:         mLock,
		cacheLatestTx: make(map[string]*model.HistoryLog),
	}
}

func (service *historyService) GetAccountLastBalance(from model.Account) uint {
	service.mLock.RLock(from.AccountNumber)
	defer service.mLock.RUnlock(from.AccountNumber)
	txHistory := service.getUserLatestTransaction(from)
	return txHistory.CurrentBalance
}

func (service *historyService) GetAccountTxHistories(from model.Account) []model.HistoryLog {
	service.mLock.RLock(from.AccountNumber)
	defer service.mLock.RUnlock(from.AccountNumber)
	return service.repo.GetAccountTxHistories(from.ID)
}

func (service *historyService) TransferBalance(from model.Account, to model.Account, nominal uint) {
	referenceID := uuid.New()

	transactionSubtract := service.SubtractBalance(from, to, nominal)
	transactionSubtract.ReferenceID = referenceID

	transactionAdd := service.AddBalance(from, to, nominal)
	transactionAdd.ReferenceID = referenceID

	transactionFee := service.SubtractBalance(from, model.Account{}, uint(transferFee))
	transactionFee.ReferenceID = referenceID

	service.ProcessTransaction()
}

func (service *historyService) ProcessTransaction() {
	if len(service.dataTransaction) > 0 {
		service.mLock.Lock("SERIAL_LOCK")         // acquire lock for bulk data transaction
		defer service.mLock.Unlock("SERIAL_LOCK") // defer release lock for bulk data transaction

		service.repo.InsertMany(service.dataTransaction)
		service.dataTransaction = service.dataTransaction[:0]
		service.cacheLatestTx = make(map[string]*model.HistoryLog)
	}
}

func (service *historyService) AddBalance(from model.Account, to model.Account, nominal uint) *model.HistoryLog {
	service.mLock.Lock(to.AccountNumber)
	defer service.mLock.Unlock(to.AccountNumber)
	var latestBalance uint
	getLatestTransaction := service.getUserLatestTransaction(to)
	latestBalance = 0
	if (*getLatestTransaction != model.HistoryLog{}) {
		latestBalance = getLatestTransaction.CurrentBalance
	}
	balanceAfterTx := latestBalance + nominal

	newTransaction := &model.HistoryLog{
		SenderID:       from.ID,
		OwnerID:        to.ID,
		LastBalance:    latestBalance,
		Balance:        int(nominal),
		CurrentBalance: balanceAfterTx,
		Tag:            "BALANCE_IN",
		CreatedAt:      time.Now(),
	}
	service.appendTransaction(newTransaction)
	service.cacheLatestTx[to.AccountNumber] = newTransaction
	return newTransaction
}

func (service *historyService) SubtractBalance(from model.Account, to model.Account, nominal uint) *model.HistoryLog {
	service.mLock.Lock(from.AccountNumber)
	defer service.mLock.Unlock(from.AccountNumber)
	var latestBalance uint

	getLatestTransaction := service.getUserLatestTransaction(from)
	latestBalance = getLatestTransaction.CurrentBalance
	balanceAfterTx := latestBalance - nominal

	newTransaction := &model.HistoryLog{
		OwnerID:        from.ID,
		RecipientID:    to.ID,
		LastBalance:    latestBalance,
		Balance:        int(nominal) * -1,
		CurrentBalance: balanceAfterTx,
		Tag:            "BALANCE_OUT",
		CreatedAt:      time.Now(),
	}
	service.appendTransaction(newTransaction)
	service.cacheLatestTx[from.AccountNumber] = newTransaction
	return newTransaction
}

func (service *historyService) appendTransaction(data *model.HistoryLog) {
	service.mLock.Lock("SERIAL_LOCK")         // acquire lock for bulk data transaction
	defer service.mLock.Unlock("SERIAL_LOCK") // defer release lock for bulk data transaction
	service.dataTransaction = append(service.dataTransaction, data)
}

func (service *historyService) getUserLatestTransaction(acc model.Account) *model.HistoryLog {
	latestTx, ok := service.cacheLatestTx[acc.AccountNumber]
	if !ok {
		latestTransaction := (service.repo.GetUserLatestTransaction(acc.ID))
		latestTx = &latestTransaction
		service.cacheLatestTx[acc.AccountNumber] = latestTx
	}
	return latestTx
}
