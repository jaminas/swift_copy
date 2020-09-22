package models

import (
	"time"
)

/**
CREATE TABLE marketplace.launch (
`event_date` Date DEFAULT toDate(event_time),
`event_time` DateTime,
`install_id` Int64,
`launch_id` UInt32,
`app_campaign_id` Int64,
`app_campaign_stream_id` Int64,
`app_id` String,
`advertising_id` String,
`device_type` UInt8,
`device_os` UInt8,
`device_model` String,
`device_country_code` FixedString(2),
`ip_country_code` FixedString(2),
`ip` UInt32,
`source_type` Enum8('none' = 0, 'push' = 1, 'adv' = 2),
`source_id` UInt32,
`sub_id_1` String,
`sub_id_2` String,
`sub_id_3` String,
`sub_id_4` String,
`sub_id_5` String,
`batch_index` UInt32
) ENGINE = MergeTree(event_date, (app_id, event_time), 8192)


0) ID
1) event_time
2) install_id
3) advertising_id
4) device_model
5) device_platform
6) device_version
7) device_serial
8) device_country_code
9) ip
10) ip_country_code
11) app_campaign_id
13) app_campaign_stream_id
12) source_type - источник запуска ( enum {“none”, “push”, “adv”} )
13) source_id - ID источника ( например пуша, по которому открыли прилу )
14) sub_id_1 ... - сабы для фильтрации
*/

type Launch struct {
	Event_date             time.Time `json:"event_date"`
	Event_time             time.Time `json:"event_time"`
	Launch_id              uint64    `json:"launch_id"`
	Install_id             int64     `json:"install_id"`
	App_campaign_id        int64     `json:"app_campaign_id"`
	App_campaign_stream_id int64     `json:"app_campaign_stream_id"`
	App_id                 string    `json:"app_id"`
	Advertising_id         string    `json:"advertising_id"`
	Device_type            uint8     `json:"device_type"`
	Device_os              uint8     `json:"device_os"`
	Device_model           string    `json:"device_model"`
	Device_country_code    string    `json:"device_country_code"`
	Ip_country_code        string    `json:"ip_country_code"`
	Ip                     uint32    `json:"ip"`
	Source_type            uint8     `json:"source_type"`
	Source_id              uint32    `json:"source_id"`

	Sid1 string `json:"sub_id_1"`
	Sid2 string `json:"sub_id_2"`
	Sid3 string `json:"sub_id_3"`
	Sid4 string `json:"sub_id_4"`
	Sid5 string `json:"sub_id_5"`
}

func NewLaunch(install *Install, request *RequestLaunch, client *Client) *Launch {

	launch := Launch{
		Event_date:             time.Now(),
		Event_time:             time.Now(),
		Launch_id:              uint64(time.Now().UnixNano()),
		Install_id:             install.Id,
		App_campaign_id:        request.CampaignId,
		App_campaign_stream_id: install.App_campaign_stream_id,
		App_id:                 install.App_id,
		Advertising_id:         install.AdvertisingId,
		Device_type:            request.DeviceType,
		Device_os:              request.DeviceOs,
		Device_model:           request.DeviceModel,
		Device_country_code:    "", // @todo
		Ip_country_code:        client.Cc2,
		Ip:                     client.IpLong,
		Source_type:            0, // @todo
		Source_id:              0, // @todo
		Sid1:                   install.Sid1,
		Sid2:                   install.Sid2,
		Sid3:                   install.Sid3,
		Sid4:                   install.Sid4,
		Sid5:                   install.Sid5,
	}

	return &launch
}
