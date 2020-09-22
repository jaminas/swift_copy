package service

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"models"
)

type ClickhouseService struct {
	dsn                string
	commit_rows        int
	accamulator_launch map[int][]models.Launch
	accamulator_view   map[int][]models.EventView
	batch_index        int
	logger             *Logger
}

/**
 * Конструктор
 */
func NewClickhouseService(logger *Logger, dsn string, commit_rows int) *ClickhouseService {

	batch_index := 0

	//var arl []models.Launch
	al := make(map[int][]models.Launch)
	al[batch_index] = []models.Launch{}

	//var arv []models.EventView
	av := make(map[int][]models.EventView)
	av[batch_index] = []models.EventView{}

	return &ClickhouseService{
		dsn:                dsn,
		commit_rows:        commit_rows,
		accamulator_launch: al,
		accamulator_view:   av,
		batch_index:        batch_index,
		logger:             logger,
	}
}

/**
 * Метод добавляет эвент в аккамулятор
 */
func (this *ClickhouseService) AddLaunch(event *models.Launch) {

	this.accamulator_launch[this.batch_index] = append(this.accamulator_launch[this.batch_index], *event)
	if len(this.accamulator_launch[this.batch_index]) >= this.commit_rows {
		go this.Commit()
	}
}

/**
 * Метод добавляет эвент в аккамулятор
 */
func (this *ClickhouseService) AddView(event *models.EventView) {

	this.accamulator_view[this.batch_index] = append(this.accamulator_view[this.batch_index], *event)
	if len(this.accamulator_view[this.batch_index]) >= this.commit_rows {
		go this.Commit()
	}
}

func (this *ClickhouseService) Commit() {

	this.logger.Info("Start to commit events to clickhouse")

	current_batch_index := this.batch_index
	this.batch_index++
	this.accamulator_launch[this.batch_index] = []models.Launch{}
	this.accamulator_view[this.batch_index] = []models.EventView{}

	go this.commit_launch(current_batch_index)
	go this.commit_view(current_batch_index)

	this.logger.Info("End to commit events to clickhouse")
}

/**
 * Метод коммитит эвенты к кликхаусу
 */
func (this *ClickhouseService) commit_launch(current_batch_index int) {

	if len(this.accamulator_launch[current_batch_index]) > 0 {

		connect, err := sql.Open("clickhouse", this.dsn)
		if err != nil {
			this.logger.Warn("Clickhouse err : " + fmt.Sprint(err))
			return
		}
		defer connect.Close()

		if err := connect.Ping(); err != nil {
			if exception, ok := err.(*clickhouse.Exception); ok {
				this.logger.Warn("Clickhouse err : " + fmt.Sprintf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
			} else {
				this.logger.Warn("Clickhouse err : " + fmt.Sprintf("%v", err))
			}
			return
		}

		tx, err := connect.Begin()
		if err != nil {
			this.logger.Warn("Clickhouse connect.Begin err : " + fmt.Sprint(err))
			return
		}

		stmt_launch, err := tx.Prepare("INSERT INTO marketplace.launch (event_date, event_time, launch_id, install_id, app_campaign_id,	app_campaign_stream_id,	app_id,	advertising_id,	device_type, device_os,	device_model, device_country_code, ip_country_code,	ip,	source_type, source_id,	sub_id_1, sub_id_2,	sub_id_3, sub_id_4,	sub_id_5, batch_index) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			this.logger.Warn("Clickhouse prepare err : " + fmt.Sprint(err))
			return
		}
		defer stmt_launch.Close()

		for _, event := range this.accamulator_launch[current_batch_index] {
			if _, err := stmt_launch.Exec(
				event.Event_time,
				event.Event_time,
				event.Launch_id,
				event.Install_id,
				event.App_campaign_id,
				event.App_campaign_stream_id,
				event.App_id,
				event.Advertising_id,
				event.Device_type,
				event.Device_os,
				event.Device_model,
				event.Device_country_code,
				event.Ip_country_code,
				event.Ip,
				event.Source_type,
				event.Source_id,
				event.Sid1,
				event.Sid2,
				event.Sid3,
				event.Sid4,
				event.Sid5,
				current_batch_index,
			); err != nil {
				this.logger.Warn("Clickhouse err : " + fmt.Sprint(err))
			}
		}

		if err := tx.Commit(); err != nil {
			this.logger.Warn("Clickhouse commit err : " + fmt.Sprint(err))

		} else {
			// чистим предыдущий массив, чтобы GC схватил эти значения
			// @todo проверить настколько качественно GC освобождает память при удалении ключа в такиз ситуациях
			delete(this.accamulator_launch, current_batch_index)
		}
	}

}

/**
 * Метод коммитит эвенты к кликхаусу
 */
func (this *ClickhouseService) commit_view(current_batch_index int) {

	if len(this.accamulator_view[current_batch_index]) > 0 {

		connect, err := sql.Open("clickhouse", this.dsn)
		if err != nil {
			this.logger.Warn("Clickhouse err : " + fmt.Sprint(err))
			return
		}
		defer connect.Close()

		if err := connect.Ping(); err != nil {
			if exception, ok := err.(*clickhouse.Exception); ok {
				this.logger.Warn("Clickhouse err : " + fmt.Sprintf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
			} else {
				this.logger.Warn("Clickhouse err : " + fmt.Sprintf("%v", err))
			}
			return
		}

		tx, err := connect.Begin()
		if err != nil {
			this.logger.Warn("Clickhouse connect.Begin err : " + fmt.Sprint(err))
			return
		}

		stmt_view, err := tx.Prepare("INSERT INTO marketplace.view (event_date, event_time, launch_id, install_id, app_campaign_id,	app_campaign_stream_id,	app_id,	advertising_id,	device_type, device_os,	device_model, user_agent, device_country_code, ip_country_code,	ip,	source_type, source_id,	sub_id_1, sub_id_2,	sub_id_3, sub_id_4,	sub_id_5, batch_index) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			this.logger.Warn("Clickhouse prepare err : " + fmt.Sprint(err))
			return
		}
		defer stmt_view.Close()

		for _, event := range this.accamulator_view[current_batch_index] {
			if _, err := stmt_view.Exec(
				event.EventTime,
				event.EventTime,
				event.LaunchId,
				event.InstallId,
				event.AppCampaignId,
				event.AppCampaignStreamId,
				event.AppId,
				event.AdvertisingId,
				event.DeviceType,
				event.DeviceOs,
				event.DeviceModel,
				event.UserAgent,
				event.DeviceCountryCode,
				event.IpCountryCode,
				event.Ip,
				event.SourceType,
				event.SourceId,
				event.Sid1,
				event.Sid2,
				event.Sid3,
				event.Sid4,
				event.Sid5,
				current_batch_index,
			); err != nil {
				this.logger.Warn("Clickhouse err : " + fmt.Sprint(err))
			}
		}

		if err := tx.Commit(); err != nil {
			this.logger.Warn("Clickhouse commit err : " + fmt.Sprint(err))

			// @todo сохранить массив в файл и очистить батч

		} else {
			// чистим предыдущий массив, чтобы GC схватил эти значения
			// @todo проверить как GC освобождает память при удалении ключа
			delete(this.accamulator_view, current_batch_index)
		}
	}

}
