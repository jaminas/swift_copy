package models

import (
	"time"
)

type EventView struct {
	EventDate           time.Time `json:"event_date"`
	EventTime           time.Time `json:"event_time"`
	InstallId           int64     `json:"install_id"`
	LaunchId            uint64    `json:"launch_id"`
	ViewId              uint64    `json:"view_id"`
	AppCampaignId       int64     `json:"app_campaign_id"`
	AppCampaignStreamId int64     `json:"app_campaign_stream_id"`
	AppId               string    `json:"app_id"`
	AdvertisingId       string    `json:"advertising_id"`
	DeviceType          uint8     `json:"device_type"`
	DeviceOs            uint8     `json:"device_os"`
	DeviceModel         string    `json:"device_model"`
	UserAgent           string    `json:"user_agent"`
	DeviceCountryCode   string    `json:"device_country_code"`
	IpCountryCode       string    `json:"ip_country_code"`
	Ip                  uint32    `json:"ip"`
	SourceType          uint8     `json:"source_type"`
	SourceId            uint32    `json:"source_id"`

	Sid1 string `json:"sub_id_1"`
	Sid2 string `json:"sub_id_2"`
	Sid3 string `json:"sub_id_3"`
	Sid4 string `json:"sub_id_4"`
	Sid5 string `json:"sub_id_5"`
}

func NewEventView(install *Install, request *RequestView, client *Client) *EventView {

	view := EventView{
		EventDate:           time.Now(),
		EventTime:           time.Now(),
		InstallId:           install.Id,
		LaunchId:            uint64(request.LaunchId),
		ViewId:              uint64(time.Now().UnixNano()),
		AppCampaignId:       install.App_campaign_id, //install.CampaignId, // @todo
		AppCampaignStreamId: install.App_campaign_stream_id,
		AppId:               install.App_id,
		AdvertisingId:       install.AdvertisingId,
		DeviceType:          request.DeviceType,
		DeviceOs:            request.DeviceOs,
		DeviceModel:         request.DeviceModel,
		UserAgent:           request.UserAgent,
		DeviceCountryCode:   "", // @todo
		IpCountryCode:       client.Cc2,
		Ip:                  client.IpLong,
		SourceType:          0, // @todo
		SourceId:            0, // @todo
		Sid1:                install.Sid1,
		Sid2:                install.Sid2,
		Sid3:                install.Sid3,
		Sid4:                install.Sid4,
		Sid5:                install.Sid5,
	}

	return &view
}
