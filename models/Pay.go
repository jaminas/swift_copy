package models

type Pay struct {
	AppInstallId             int64   `json:"app_install_id"`
	AppInstallPayConditionId int64   `json:"app_install_pay_condition_id"`
	Cost                     float32 `json:"cost"`
	Commission               float32 `json:"commission"`
	CurrencyCode             string  `json:"currency_code"`
}

func NewPay(install *Install, pay_condition *PayCondition) *Pay {

	pay := Pay{
		AppInstallId:             install.Id,
		AppInstallPayConditionId: pay_condition.Id,
		Cost:                     pay_condition.Cost,
		Commission:               pay_condition.Commission,
		CurrencyCode:             pay_condition.CurrencyCode,
	}

	return &pay
}
