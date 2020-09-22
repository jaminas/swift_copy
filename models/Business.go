package models

type Business struct {
	InstallId string          `json:"install_id"`
	Request   *RequestLaunch  `json:"request"`
	Client    *Client         `json:"client"`
	Campaign  *Campaign       `json:"campaign"`
	Stream    *CampaignStream `json:"stream"`
	Install   *Install        `json:"install"`
	Pay       *Pay            `json:"pay"`
	Launch    *Launch         `json:"launch"`
}
