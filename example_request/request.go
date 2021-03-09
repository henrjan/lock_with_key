package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	accountListFrom = []string{
		"104938043525",
		"442897384221",
		"037733806596",
		"465179597739",
		"036073356959",
		"980055396477",
		"769369146010",
		"948322084003",
		"368834496660",
	}
	accountDestination = "154687370751"
	nominal            = 10000
)

func main() {
	var wg sync.WaitGroup

	for _, v := range accountListFrom {
		wg.Add(1)
		go func(accountFrom string) {
			defer wg.Done()

			postBody, _ := json.Marshal(map[string]interface{}{
				"account_from": accountFrom,
				"account_to":   accountDestination,
				"nominal":      uint(nominal),
			})

			requestBody := bytes.NewBuffer(postBody)

			fmt.Printf("nominal value : %v\n", requestBody)
			resp, err := http.Post(
				"http://127.0.0.1:8080/transaction/transfer",
				"application/json",
				requestBody,
			)
			defer resp.Body.Close()
			if err != nil {
				panic(err)
			}
		}(v)
	}
	wg.Wait()
}
