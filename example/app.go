package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/henrjan/lock_with_key/example/general"
	"github.com/henrjan/lock_with_key/example/model"
	"github.com/henrjan/lock_with_key/example/repository"
	"github.com/henrjan/lock_with_key/example/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DB_NAME = "data.db"
)

var (
	multiLock = general.NewMultipleLock()

	db             = initDB()
	accountRepo    = repository.NewAccountRepository(db)
	historyRepo    = repository.NewHistoryRepository(db)
	accountService = service.NewAccountService(accountRepo)
	historyService = service.NewHistoryService(historyRepo, multiLock)
)

func main() {

	app := fiber.New(fiber.Config{})
	insertDefaultBalanceToAllAccount()

	app.Get("/transaction/last_balance", func(c *fiber.Ctx) error {
		accNumber := c.Query("account_number")

		getAccount := accountService.GetAccountByNumber(accNumber)
		balance := historyService.GetAccountLastBalance(getAccount)

		result := fiber.Map{
			"result": balance,
			"error":  nil,
		}

		return c.JSON(result)
	})

	app.Get("/transaction/history_log", func(c *fiber.Ctx) error {
		accNumber := c.Query("account_number")

		getAccount := accountService.GetAccountByNumber(accNumber)
		txHistories := historyService.GetAccountTxHistories(getAccount)

		result := fiber.Map{
			"result": txHistories,
			"error":  nil,
		}

		return c.JSON(result)
	})

	app.Post("/transaction/transfer", func(c *fiber.Ctx) error {
		body := struct {
			AccountFrom string `json:"account_from"`
			AccountTo   string `json:"account_to"`
			Nominal     uint   `json:"nominal"`
		}{}

		if err := json.Unmarshal(c.Body(), &body); err != nil {
			result := fiber.Map{
				"result": nil,
				"error":  err.Error(),
			}
			fmt.Printf("errors : %v", err)
			return c.JSON(result)
		}

		accFrom := accountService.GetAccountByNumber(body.AccountFrom)
		accTo := accountService.GetAccountByNumber(body.AccountTo)

		timeNow := time.Now()
		historyService.TransferBalance(accFrom, accTo, body.Nominal)
		duration := time.Since(timeNow)
		fmt.Printf("duration : %d\n", duration)

		result := fiber.Map{
			"result": "success",
			"error":  nil,
		}

		return c.JSON(result)
	})

	app.Listen(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func insertDefaultBalanceToAllAccount() {
	balance := 100000
	accountList := accountService.GetAllAccount()

	for _, v := range accountList {
		historyService.AddBalance(model.Account{}, v, uint(balance), uuid.UUID{})
	}
}
