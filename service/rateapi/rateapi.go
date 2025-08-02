package rateapi

import (
	"encoding/json"
	"fmt"
	"github.com/albegonzalezp/ratesupdater/models"
	"io"
	"net/http"
	"os"
)

type Service struct {
	Url    string
	ApiKey string
	Client *http.Client
}

func NewService(url string, apiKey string, client *http.Client) *Service {
	return &Service{Url: url, ApiKey: apiKey, Client: client}
}

func (s *Service) GetRates(baseCode string) (rates models.ExchangeRateResponse, err error) {

	resp, err := s.Client.Get(s.Url + os.Getenv("EXCHANGE_API_KEY") + "/latest/" + baseCode)
	if err != nil {
		return rates, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return rates, err
	}

	err = json.Unmarshal(bs, &rates)
	if err != nil {
		return rates, err
	}

	if rates.Result != "success" {
		return rates, fmt.Errorf(rates.Result)
	}

	return rates, nil

}
