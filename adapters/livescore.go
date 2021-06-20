package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yudgnahk/euro21/dtos"
)

const (
	Host      = "https://prod-public-api.livescore.com/v1/api/react"
	TablePath = "leagueTable/soccer/euro-2020"
)

func newGetRequest(url string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, url, nil)
}

func setHeaders(req *http.Request) {
	req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Add("Referer", "https://www.livescore.com/")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
}

func execute(req *http.Request, response interface{}) error {
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}

func GetTables() (*dtos.TableData, error) {
	request, _ := newGetRequest(fmt.Sprintf("%v/%v", Host, TablePath))
	setHeaders(request)

	var response dtos.TableData

	err := execute(request, &response)
	return &response, err
}
