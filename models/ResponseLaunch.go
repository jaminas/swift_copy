package models

type ResponseLaunch struct {
	InstallId    string    `json:"install_id"`
	AllowWebview bool      `json:"allow_webview"`
	Business     *Business `json:"business"`
}
