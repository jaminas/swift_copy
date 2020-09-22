package models

type Campaign struct {
	Id                    int64            `json:"id"`
	App_id                string           `json:"app_id"`
	Integrated_partner_id string           `json:"integrated_partner_id"`
	Streams               []CampaignStream `json:"streams"`
	Modules               []AppModule      `json:"modules"`
}

type CampaignStream struct {
	Id        int64    `json:"id"`
	Url       string   `json:"url"`
	Weight    uint8    `json:"weight"`
	Countries []string `json:"countries"`
}

type AppModule struct {
	Id       int64  `json:"id"`
	settings string `json:"settings"`
}
