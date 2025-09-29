package models

import "strconv"

type P2pRateRequest struct {
	Page      int    `json:"page"`
	Rows      int    `json:"rows"`
	Asset     string `json:"asset"`
	Fiat      string `json:"fiat"`
	TradeType string `json:"tradeType"`
	// PayTypes  []string `json:"payTypes"`
}

// {"page":1,"rows":10,"asset":"USDT","fiat":"VES","tradeType":"SELL","payTypes":[]}

type P2pRateResponse struct {
	Code string `json:"code"`
	Data []Data `json:"data"`
}

type Data struct {
	Adv Adv `json:"adv"`
}

type Adv struct {
	AdvNo     string `json:"advNo"`
	Asset     string `json:"asset"`
	Price     string `json:"price"`
	TradeType string `json:"tradeType"`
	Clasify   string `json:"clasify"`
	FiatUnit  string `json:"fiatUnit"`
}

// Convert to USD/VES rate
func (r *P2pRateResponse) ToRateDB() (rates Rate) {

	if len(r.Data) == 0 {
		return Rate{}
	}

	for _, d := range r.Data {
		if d.Adv.Asset == "USDT" && d.Adv.FiatUnit == "VES" {

			price, err := strconv.ParseFloat(d.Adv.Price, 64)
			if err != nil {
				return Rate{}
			}

			return Rate{
				BaseCode: "USD",
				Code:     "VES",
				Rate:     price,
			}
		}
	}

	return Rate{}
}
