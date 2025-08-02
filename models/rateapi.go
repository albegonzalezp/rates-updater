package models

// ExchangeRateResponse is the response from exhange rate api
type ExchangeRateResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int                `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int                `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

func (r *ExchangeRateResponse) ToRateDB() (dbRates []Rate) {
	for code, rate := range r.ConversionRates {
		r := Rate{
			BaseCode: r.BaseCode,
			Rate:     rate,
			Code:     code,
		}

		dbRates = append(dbRates, r)
	}
	return
}
