package models

type ResponseView struct {
	RedirectUrl string       `json:"redirect_url"`
	Request     *RequestView `json:"request"`
	Client      *Client      `json:"client"`
	Install     *Install     `json:"install"`
	EventView   *EventView   `json:"event_view"`
}
