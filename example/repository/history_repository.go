package repository

import (
	"github.com/henrjan/lock_with_key/example/model"
	"gorm.io/gorm"
)

type HistoryRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(dbClient *gorm.DB) *HistoryRepository {
	dbClient.AutoMigrate(&model.HistoryLog{})
	dbClient.Where("1 = 1").Delete(&model.HistoryLog{})
	return &HistoryRepository{dbClient}
}

func (repo *HistoryRepository) InsertMany(data []model.HistoryLog) {
	repo.db.Create(&data)
}
func (repo *HistoryRepository) Insert(data model.HistoryLog) {
	repo.db.Create(&data)
}

func (repo *HistoryRepository) GetUserLatestTransaction(ownerID uint) model.HistoryLog {
	historyLogs := model.HistoryLog{}
	repo.db.Order("created_at desc").Limit(1).Where("owner_id = ?", ownerID).Find(&historyLogs)
	return historyLogs
}

func (repo *HistoryRepository) GetAccountTxHistories(ownerID uint) []model.HistoryLog {
	txHistories := make([]model.HistoryLog, 0)
	repo.db.Order("id desc").Where("owner_id = ?", ownerID).Find(&txHistories)
	return txHistories
}
