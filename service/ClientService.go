package service

import (
	"github.com/astaxie/beego/context"
	"misc"
	"models"
	"strings"
)

type ClientService struct {
	logger *Logger
}

/**
 * Конструктор
 */
func NewClientService(logger *Logger) *ClientService {
	return &ClientService{logger: logger}
}

const (
	DEVICE_TYPE_UNKNOWN uint8 = 0
	DEVICE_TYPE_MOBILE  uint8 = 1
	DEVICE_TYPE_TABLET  uint8 = 2
	DEVICE_TYPE_PC      uint8 = 3

	DEVICE_OS_UNKNOWN       uint16 = 0
	DEVICE_OS_ANDROID       uint16 = 1
	DEVICE_OS_IOS           uint16 = 2
	DEVICE_OS_BLACKBERRY    uint16 = 3
	DEVICE_OS_SYMBOS        uint16 = 4
	DEVICE_OS_WINDOWS_PHONE uint16 = 5
	DEVICE_OS_NOKIA         uint16 = 6
	DEVICE_OS_WINDOWS       uint16 = 7
	DEVICE_OS_OSX           uint16 = 8
	DEVICE_OS_OS2           uint16 = 9
	DEVICE_OS_SUNOS         uint16 = 10
	DEVICE_OS_CHROME_OS     uint16 = 11
	DEVICE_OS_FREEBSD       uint16 = 12
	DEVICE_OS_NETBSD        uint16 = 13
	DEVICE_OS_OPENBSD       uint16 = 14
	DEVICE_OS_OPENSOLARIS   uint16 = 15
	DEVICE_OS_BEOS          uint16 = 16
	DEVICE_OS_LINUX         uint16 = 17
)

func (this *ClientService) Parse(ctx *context.Context) *models.Client {

	country_code := ctx.Input.Query("cc2") // костыль для передачи страны через запрос
	if country_code == "" {
		country_code = ctx.Input.Header("country_code")
	}

	client := models.Client{
		Cc2:    country_code,
		Ip:     this.getClientIp(ctx),
		IpLong: misc.Ip2uint32(this.getClientIp(ctx)),
	}

	return &client
}

func (this *ClientService) getClientIp(ctx *context.Context) string {
	s := strings.Split(ctx.Input.Header("X-Real-IP"), ":")
	return s[0]
}
