package service

import (
	"database/sql"
	"fmt"
	"models"
	_ "time"
)

type PayService struct {
	db              *sql.DB
	logger          *Logger
	condition_store map[string]models.PayCondition // хранилице всех условий оплаты
	accamulator     map[int][]models.Pay           // аккамулятор для платежей на запись
	batch_index     int
	commit_rows     int
}

func NewPayService(logger *Logger, db *sql.DB, commit_rows int) *PayService {

	batch_index := 0
	al := make(map[int][]models.Pay)
	al[batch_index] = []models.Pay{}

	return &PayService{
		db:              db,
		logger:          logger,
		condition_store: make(map[string]models.PayCondition),
		accamulator:     al,
		batch_index:     batch_index,
		commit_rows:     commit_rows,
	}
}

/**
 * Метод добавляет в аккамулятор платеж
 */
func (this *PayService) Add(pay *models.Pay) {

	this.accamulator[this.batch_index] = append(this.accamulator[this.batch_index], *pay)
	if len(this.accamulator[this.batch_index]) >= this.commit_rows {
		go this.Commit()
	}
}

/**
 * Метод коммитит все пейменты
 */
func (this *PayService) Commit() {

	this.logger.Info("Start to commit pays to database. Size: " + fmt.Sprint(len(this.accamulator[this.batch_index])))

	current_batch_index := this.batch_index

	new_batch_index := this.batch_index + 1
	this.accamulator[new_batch_index] = []models.Pay{}
	this.batch_index = new_batch_index

	statement := `INSERT INTO tbl_app_install_pay (app_install_id, app_install_pay_condition_id, cost, commission, currency_code) VALUES ($1, $2, $3, $4, $5)`
	for _, pay := range this.accamulator[current_batch_index] {

		_, err := this.db.Exec(statement,
			pay.AppInstallId,
			pay.AppInstallPayConditionId,
			pay.Cost,
			pay.Commission,
			pay.CurrencyCode,
		)
		if err != nil {
			this.logger.Warn(fmt.Sprintf("Pay save error: %v", err))
		}
	}
	delete(this.accamulator, current_batch_index)

	this.logger.Info("End to commit pays to database")
}

/**
 * Метод возвращает ID условия оплаты
 */
func (this *PayService) GetPayCondition(client *models.Client, campaign *models.Campaign) *models.PayCondition {

	if val, ok := this.condition_store[campaign.App_id+"_"+client.Cc2]; ok {
		return &val
	} else {
		if val, ok := this.condition_store[campaign.App_id+"_"]; ok { // базовые условия
			return &val
		}
	}
	return nil
}

/**
 * Метод производит выборку всех условий оплаты и сохраняет в стор
 */
func (this *PayService) Load2Store() {

	this.logger.Info("Start to load conditions to store")

	cst := make(map[string]models.PayCondition)

	sql := `SELECT id, app_id, COALESCE(country_code::text,''), cost, commission, currency_code FROM tbl_app_install_pay_condition WHERE deleted_at IS NULL`
	if rows, err := this.db.Query(sql); err == nil {

		defer rows.Close()
		for rows.Next() {

			var pcond models.PayCondition

			if err := rows.Scan(&pcond.Id, &pcond.AppId, &pcond.CountryCode, &pcond.Cost, &pcond.Commission, &pcond.CurrencyCode); err == nil {
				cst[pcond.AppId+"_"+pcond.CountryCode] = pcond
			} else {
				this.logger.Warn(fmt.Sprintf("condition scan error: %v", err))
			}
		}
		if len(cst) > 0 { // только в случае присутсвия элементов меняем стор

			this.condition_store = cst
		} else {

			this.logger.Warn("conditions are empty")
		}
	} else {

		this.logger.Warn(fmt.Sprintf("condition query error: %v", err))
	}

	this.logger.Info("End to load conditions to store. Size: " + fmt.Sprint(len(cst)))

}
