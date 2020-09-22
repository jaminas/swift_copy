package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Install struct {
	Id                             int64  `json:"id"`
	App_id                         string `json:"app_id"`
	Push_module_id                 string `json:"push_module_id"`
	Push_token                     string `json:"push_token"`
	Api_domain_id                  int64  `json:"domain_id"`
	AdvertisingId                  string `json:"adv_id"`
	Integrated_partner_id          string `json:"ip_id"`
	Integrated_partner_custom_data string `json:"ipcd"`
	Allow_webview                  bool   `json:"aw"`
	Sid1                           string `json:"sid1"`
	Sid2                           string `json:"sid2"`
	Sid3                           string `json:"sid3"`
	Sid4                           string `json:"sid4"`
	Sid5                           string `json:"sid5"`

	App_campaign_id                          int64   `json:"cmp_id"`
	App_campaign_stream_id                   int64   `json:"str_id"`
	Webview_url                              string  `json:"wurl"`
	Integrated_partner_install_cost          float32 `json:"ipic"`
	Integrated_partner_install_cost_currency string  `json:"ipicc"`
}

func NewInstall(request *RequestLaunch, client *Client, campaign *Campaign, stream *CampaignStream, allow_webview bool) *Install {

	ipcd := "{}"

	// @todo вынести это в реквест
	if request.Ipcd != nil {
		if json_ipcd, err := json.Marshal(request.Ipcd); err == nil {
			ipcd = string(json_ipcd)
		} else {
			fmt.Sprintln("ipcd json marshaling error: ", err)
		}
	}

	install := Install{
		Id:                                       time.Now().UnixNano(),
		App_id:                                   campaign.App_id,
		Push_module_id:                           request.PushModuleId,
		Push_token:                               request.PushToken,
		Api_domain_id:                            request.DomainId,
		AdvertisingId:                            request.AdvertisingId,
		Integrated_partner_id:                    campaign.Integrated_partner_id,
		Integrated_partner_custom_data:           ipcd,
		Allow_webview:                            allow_webview,
		Sid1:                                     request.Sid1,
		Sid2:                                     request.Sid2,
		Sid3:                                     request.Sid3,
		Sid4:                                     request.Sid4,
		Sid5:                                     request.Sid5,
		App_campaign_id:                          campaign.Id,
		App_campaign_stream_id:                   stream.Id,
		Webview_url:                              stream.Url,
		Integrated_partner_install_cost:          request.Ipic,
		Integrated_partner_install_cost_currency: request.Ipicc,
	}

	return &install
}
