package models

type RequestLaunch struct {

	// данные из хвоста
	CampaignId int64 `json:"campaign_id"`
	DomainId   int64 `json:"domain_id"`

	Sid1 string `json:"sub_id_1"`
	Sid2 string `json:"sub_id_2"`
	Sid3 string `json:"sub_id_3"`
	Sid4 string `json:"sub_id_4"`
	Sid5 string `json:"sub_id_5"`

	AdvertisingId string `json:"advertising_id"`
	PushModuleId  string `json:"push_module_id"`
	PushToken     string `json:"push_token"`
	//Install       string  `json:"install"` // hashid инстала
	Ipcd  map[string]interface{} `json:"ipcd"`  //integrated_partner_custom_data
	Ipic  float32                `json:"ipic"`  //Integrated_partner_install_cost
	Ipicc string                 `json:"ipicc"` //Integrated_partner_install_cost_currency

	Install       string   `json:"install"`        // хеш install
	InstallObject *Install `json:"install_object"` // объект install из кеша

	DeviceType  uint8  `json:"device_type"`
	DeviceModel string `json:"device_model"`
	DeviceOs    uint8  `json:"device_os"`
}
