package module

import (
	//"models"
	"service"
	//"github.com/astaxie/beego/context"
	//"strings"
)

type GoogleBotService struct {
	logger *service.Logger
}

func NewGoogleBotService(logger *service.Logger) *GoogleBotService {
	return &GoogleBotService{logger: logger}
}

const (
	DEVICE_TYPE_UNKNOWN uint8 = 0
	DEVICE_TYPE_MOBILE  uint8 = 1
	DEVICE_TYPE_TABLET  uint8 = 2
	DEVICE_TYPE_PC      uint8 = 3
)
