package models

type RequestView struct {

	// данные из хвоста
	InstallObject *Install `json:"install_object"`
	LaunchId      int64    `json:"launch_id"`

	UserAgent string `json:"user_agent"`

	DeviceType  uint8  `json:"device_type"`
	DeviceModel string `json:"device_model"`
	DeviceOs    uint8  `json:"device_os"`
}
