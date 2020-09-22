package service

import (
	"database/sql"
	"fmt"
	"math/rand"
	"misc"
	"models"
	"time"
	_ "time"
)

type CampaignService struct {
	db     *sql.DB
	logger *Logger
	store  map[int64]models.Campaign
}

/**
 * Конструктор
 */
func NewCampaignService(logger *Logger, db *sql.DB) *CampaignService {
	return &CampaignService{db: db, logger: logger}
}

/**
 * Метод возвращает стрим
 */
func (this *CampaignService) GetCampaign(id int64) *models.Campaign {

	campaign := models.Campaign{}
	if val, ok := this.store[id]; ok {
		campaign = val
	}
	return &campaign
}

/**
 * Метод находит стрим победитель из стримов кампании
 */
func (this *CampaignService) GetWinnerStream(campaign *models.Campaign, request *models.RequestLaunch, client *models.Client) *models.CampaignStream {

	stream := models.CampaignStream{}
	streams := []models.CampaignStream{}
	var length int = 0 // длина весов

	for _, item := range campaign.Streams {

		if misc.ContainsInArray(item.Countries, client.Cc2) {

			streams = append(streams, item)

			stream_length := int(item.Weight)
			if stream_length == 0 {
				stream_length = 1
			}

			length = length + stream_length
		}
	}

	if length == 0 {
		return &stream
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randnum := r1.Intn(length) + 1

	// обнуляем длину для переиспользования
	length = 0

	for _, item := range streams {

		stream_length := int(item.Weight)
		if stream_length == 0 {
			stream_length = 1
		}
		length = length + stream_length

		if randnum <= length {

			return &item
		}
	}

	return &stream
}

/**
 * Метод производит выборку всех кампаний и сохраняет в стор
 */
func (this *CampaignService) Load2Store() {

	this.logger.Info("Start to load campaigns to store")

	st := make(map[int64]models.Campaign)

	sql := `SELECT id, app_id, COALESCE(integrated_partner_id, '') FROM tbl_app_campaign`
	if rows, err := this.db.Query(sql); err == nil {

		defer rows.Close()
		for rows.Next() {

			cmp := models.Campaign{}
			if err := rows.Scan(&cmp.Id, &cmp.App_id, &cmp.Integrated_partner_id); err == nil {
				st[cmp.Id] = cmp
			} else {
				this.logger.Warn(fmt.Sprintf("campaign scan error: %v", err))
			}
		}

	} else {
		this.logger.Warn(fmt.Sprintf("campaigns query error: %v", err))
	}

	// парсим стримы кампании
	sql = `SELECT id, url, weight, available_country_code_list, app_campaign_id FROM tbl_app_campaign_stream`
	if rows, err := this.db.Query(sql); err == nil {

		defer rows.Close()
		for rows.Next() {

			str := models.CampaignStream{}
			var (
				campaign_id  int64
				country_list string
			)

			if err := rows.Scan(&str.Id, &str.Url, &str.Weight, &country_list, &campaign_id); err == nil {

				if cmp, ok := st[campaign_id]; ok {

					str.Countries = misc.Enum2Array(country_list)
					cmp.Streams = append(cmp.Streams, str)
					st[cmp.Id] = cmp
				}
			} else {
				this.logger.Warn(fmt.Sprintf("streams scan error: %v", err))
			}
		}

	} else {
		this.logger.Warn(fmt.Sprintf("streams query error: %v", err))
	}

	if len(st) > 0 { // только в случае присутсвия элементов меняем стор
		this.store = st
	} else {
		this.logger.Warn("campaigns are empty")
	}

	this.logger.Info("End to load campaigns to store. Size: " + fmt.Sprint(len(st)))

}
