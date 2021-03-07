package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/henrjan/lock_with_key/example/general"
	"github.com/henrjan/lock_with_key/example/model"
	"github.com/henrjan/lock_with_key/example/repository"
	"github.com/henrjan/lock_with_key/example/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	counter = 0

	wg sync.WaitGroup
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
		fmt.Printf("body request : %v\n", string(c.Body()))

		accFrom := accountService.GetAccountByNumber(body.AccountFrom)
		accTo := accountService.GetAccountByNumber(body.AccountTo)

		historyService.TransferBalance(accFrom, accTo, body.Nominal)

		result := fiber.Map{
			"result": "success",
			"error":  nil,
		}

		return c.JSON(result)
	})

	app.Listen(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func insertDefaultBalanceToAllAccount() {
	balance := 100000
	accountList := accountService.GetAllAccount()

	for _, v := range accountList {
		historyService.AddBalance(model.Account{}, v, uint(balance))
	}
	historyService.ProcessTransaction()
}

func transferBalance(wg *sync.WaitGroup) {
	accountList := accountService.GetAllAccount()

	n := len(accountList) - 1
	for k, v := range accountList {
		if k == n {
			break
		}
		wg.Add(1)
		go func(data model.Account) {
			defer wg.Done()
			historyService.TransferBalance(data, accountList[n], 10000)
		}(v)
	}
}
