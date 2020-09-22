package models

type PayCondition struct {
	Id           int64   `json:"id"`
	AppId        string  `json:"app_id"`
	CountryCode  string  `json:"country_code"`
	Cost         float32 `json:"cost"`
	Commission   float32 `json:"commission"`
	CurrencyCode string  `json:"currency_code"`
}
