package rateapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/albegonzalezp/ratesupdater/models"
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
		return rates, fmt.Errorf("%v", rates.Result)
	}

	return rates, nil

}

// Get USDT/VES rate
func (s *Service) GetUSDTtoVESp2pRate() (p2pResp models.P2pRateResponse, err error) {

	req := models.P2pRateRequest{
		Page:      1,
		Rows:      10,
		Asset:     "USDT",
		Fiat:      "VES",
		TradeType: "BUY",
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return p2pResp, err
	}

	resp, err := s.Client.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return p2pResp, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return p2pResp, err
	}

	err = json.Unmarshal(bs, &p2pResp)
	if err != nil {
		return p2pResp, err
	}

	return p2pResp, nil
}
