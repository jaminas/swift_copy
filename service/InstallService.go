package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"misc"
	"models"
	_ "time"
)

type InstallService struct {
	db          *sql.DB
	logger      *Logger
	cache       cache.Cache
	accamulator map[int][]models.Install // аккамулятор для инсталов на запись
	batch_index int
	commit_rows int
}

/**
 * Конструктор
 */
func NewInstallService(
	logger *Logger,
	db *sql.DB,
	commit_rows int,
) *InstallService {

	bm, err := cache.NewCache("memory", `{}`) //`{"interval":10}`
	if err != nil {
		logger.Warn(fmt.Sprintf("%v", err))
	}

	batch_index := 0
	al := make(map[int][]models.Install)
	al[batch_index] = []models.Install{}

	return &InstallService{
		db:          db,
		logger:      logger,
		cache:       bm,
		accamulator: al,
		batch_index: batch_index,
		commit_rows: commit_rows,
	}
}

/**
 * Метод возвращает инсталл
 */
func (this *InstallService) Get(id int64) *models.Install {

	install := this.getFromCache(id)
	if install.Id == 0 {

		sql := `SELECT i.id, i.app_id, COALESCE(i.push_module_id, ''), i.push_token, COALESCE(i.api_domain_id, 0), COALESCE(i.advertising_id, ''), i.integrated_partner_id, i.integrated_partner_custom_data, i.allow_webview, i.sub_id_1, i.sub_id_2, i.sub_id_3, i.sub_id_4, i.sub_id_5, COALESCE(i.app_campaign_stream_id, 0),  COALESCE(i.webview_url, ''), i.integrated_partner_install_cost, COALESCE(i.integrated_partner_install_cost_currency, ''), COALESCE(i.app_campaign_id, 0) FROM tbl_app_install i WHERE i.id = $1`

		err := this.db.QueryRow(sql, id).Scan(
			&install.Id,
			&install.App_id,
			&install.Push_module_id,
			&install.Push_token,
			&install.Api_domain_id,
			&install.AdvertisingId,
			&install.Integrated_partner_id,
			&install.Integrated_partner_custom_data,
			&install.Allow_webview,
			&install.Sid1,
			&install.Sid2,
			&install.Sid3,
			&install.Sid4,
			&install.Sid5,
			&install.App_campaign_stream_id,
			&install.Webview_url,
			&install.Integrated_partner_install_cost,
			&install.Integrated_partner_install_cost_currency,
			&install.App_campaign_id,
		)
		if err != nil {
			this.logger.Warn(fmt.Sprintf(sql+", id=%v : %v", id, err))
		}
	}

	return &install
}

/**
 * Метод добавляет в аккамулятор инсталл
 */
func (this *InstallService) Add(install *models.Install) {

	this.accamulator[this.batch_index] = append(this.accamulator[this.batch_index], *install)
	if len(this.accamulator[this.batch_index]) >= this.commit_rows {
		go this.Commit()
	}
}

/**
 * Метод коммитит все пейменты
 */
func (this *InstallService) Commit() {

	this.logger.Info("Start to commit installs to database. Size: " + fmt.Sprint(len(this.accamulator[this.batch_index])))

	current_batch_index := this.batch_index

	new_batch_index := this.batch_index + 1
	this.accamulator[new_batch_index] = []models.Install{}
	this.batch_index = new_batch_index

	statement := `INSERT INTO tbl_app_install (id, app_id, push_module_id, push_token, api_domain_id, advertising_id, integrated_partner_id, integrated_partner_custom_data, allow_webview, sub_id_1, sub_id_2, sub_id_3, sub_id_4, sub_id_5, app_campaign_stream_id, webview_url, integrated_partner_install_cost, integrated_partner_install_cost_currency) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	for _, install := range this.accamulator[current_batch_index] {
		_, err := this.db.Exec(statement,
			install.Id,
			misc.StringSql(install.App_id),
			misc.StringSql(install.Push_module_id),
			install.Push_token,
			misc.Int64Sql(install.Api_domain_id),
			misc.StringSql(install.AdvertisingId),
			misc.StringSql(install.Integrated_partner_id),
			install.Integrated_partner_custom_data,
			install.Allow_webview,
			install.Sid1,
			install.Sid2,
			install.Sid3,
			install.Sid4,
			install.Sid5,
			misc.Int64Sql(install.App_campaign_stream_id),
			misc.StringSql(install.Webview_url),
			install.Integrated_partner_install_cost,
			misc.StringSql(install.Integrated_partner_install_cost_currency),
		)
		if err != nil {
			this.logger.Warn(fmt.Sprintf("Install save error: ", err))
		}
	}
	delete(this.accamulator, current_batch_index)

	this.logger.Info("End to commit installs to database")
}

/**
 * Метод возвращает стрим из кеша
 */
func (this *InstallService) getFromCache(id int64) models.Install {

	install := models.Install{}
	cachekey := fmt.Sprintf("install:%v", id)

	if this.cache.IsExist(cachekey) {
		cachevalue := this.cache.Get(cachekey)
		cachedata := string(cachevalue.([]uint8))
		if cachedata != "" {
			b := bytes.NewBufferString(cachedata)
			if err := json.Unmarshal(b.Bytes(), &install); err != nil {
				this.logger.Warn(fmt.Sprintf("%v", err))
			}
		}
	}

	return install
}
