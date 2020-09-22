package service

import (
	"fmt"
	"github.com/speps/go-hashids"
	"models"
	"strings"
	_ "time"
)

type MacrosService struct {
	logger *Logger
	HashId *hashids.HashID
}

/**
 * Конструктор
 */
func NewMacrosService(
	logger *Logger,
	hash_id *hashids.HashID,
) *MacrosService {

	return &MacrosService{
		logger: logger,
		HashId: hash_id,
	}
}

/**
 * Метод вставляет макросы
 */
func (this *MacrosService) FillMacros(url string, install *models.Install, client *models.Client, event_view *models.EventView) string {

	url = strings.ReplaceAll(url, "{install_id}", fmt.Sprintf("%v", install.Id))
	url = strings.ReplaceAll(url, "{launch_id}", fmt.Sprintf("%v", event_view.LaunchId))
	url = strings.ReplaceAll(url, "{view_id}", fmt.Sprintf("%v", event_view.ViewId))
	url = strings.ReplaceAll(url, "{campaign_id}", fmt.Sprintf("%v", install.App_campaign_id))
	url = strings.ReplaceAll(url, "{stream_id}", fmt.Sprintf("%v", install.App_campaign_stream_id))
	url = strings.ReplaceAll(url, "{app_id}", install.App_id)
	url = strings.ReplaceAll(url, "{sub_id_1}", install.Sid1)
	url = strings.ReplaceAll(url, "{sub_id_2}", install.Sid2)
	url = strings.ReplaceAll(url, "{sub_id_3}", install.Sid3)
	url = strings.ReplaceAll(url, "{sub_id_4}", install.Sid4)
	url = strings.ReplaceAll(url, "{sub_id_5}", install.Sid5)

	url = strings.ReplaceAll(url, "{country_code}", client.Cc2)
	//url = strings.ReplaceAll(url, "{fp_1}", request_data.Fp1)
	//url = strings.ReplaceAll(url, "{fp_2}", request_data.Fp2)
	//url = strings.ReplaceAll(url, "{fp_3}", request_data.Fp3)
	//url = strings.ReplaceAll(url, "{fp_4}", request_data.Fp4)
	//url = strings.ReplaceAll(url, "{fp_5}", request_data.Fp5)

	viewkey := ""
	var numbers []int64

	var znachenie = install.Id
	var znak int64 = 1
	if znachenie < 0 {
		znak = 0
		znachenie = -znachenie
	}

	numbers = append(numbers, znak, znachenie, int64(event_view.ViewId))
	if vhash, err := this.HashId.EncodeInt64(numbers); err == nil {
		viewkey = vhash
	} else {
		this.logger.Warn(fmt.Sprintf("%v", err))
	}

	url = strings.ReplaceAll(url, "{viewkey}", viewkey)

	return url

}
