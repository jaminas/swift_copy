package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/speps/go-hashids"
	"models"
	_ "time"
)

type BusinessLayer struct {
	Logger            *Logger
	CampaignService   *CampaignService
	InstallService    *InstallService
	PayService        *PayService
	ClickhouseService *ClickhouseService
	ModuleService     *ModuleService
	HashId            *hashids.HashID
}

func NewBusinessLayer(
	Logger *Logger,
	CampaignService *CampaignService,
	InstallService *InstallService,
	PayService *PayService,
	ClickhouseService *ClickhouseService,
	ModuleService *ModuleService,
	HashId *hashids.HashID,
) *BusinessLayer {

	return &BusinessLayer{
		Logger:            Logger,
		CampaignService:   CampaignService,
		InstallService:    InstallService,
		PayService:        PayService,
		ClickhouseService: ClickhouseService,
		ModuleService:     ModuleService,
		HashId:            HashId,
	}
}

func (this *BusinessLayer) Get(request *models.RequestLaunch, client *models.Client) *models.Business {

	response := models.Business{
		Request: request,
		Client:  client,
	}

	// в случае если Hash не верный, то дальше не идем
	if request.CampaignId == 0 {
		return &response
	}

	// находим кампанию
	campaign := this.CampaignService.GetCampaign(request.CampaignId)
	response.Campaign = campaign
	if campaign.Id == 0 { // почти невероятно, но если объект не найден, то выдаем 400
		return &response
	}

	var install *models.Install

	// проверяем install_object, если его не существует , то генерим его. Если существует, находим
	if request.InstallObject == nil {

		// шаг1. Нашли победителя стримов
		stream := this.CampaignService.GetWinnerStream(campaign, request, client)
		response.Stream = stream

		// шаг2. Пройти все валидации модулей
		check_modules := this.ModuleService.Check(request, client, campaign, stream)

		// шаг3. если стрим нашелся и проверки модулей пройдены, то allow_webview=true , иначе false
		allow_webview := stream.Id != 0 && check_modules

		// создаем инстралл
		install = models.NewInstall(request, client, campaign, stream, allow_webview)
		this.InstallService.Add(install)

		// создаем pay
		pay_condition := this.PayService.GetPayCondition(client, campaign)
		if pay_condition != nil {
			pay := models.NewPay(install, pay_condition)
			this.PayService.Add(pay)
			response.Pay = pay
		}

	} else {

		// находим инсталл в кеше или базе
		install = this.InstallService.Get(request.InstallObject.Id)
		if install.Id == 0 { // инсталл может быть не найден в кеше или базе если это другая нода или не произошла синхронизация с базой

			// восстанавливаем инсталл из запроса - именно поэтому мы передаем объект инстала в запросе
			install = request.InstallObject
		}
	}
	response.Install = install

	var launch *models.Launch = models.NewLaunch(install, request, client)
	response.Launch = launch

	// записываем лаунч в КХ
	this.ClickhouseService.AddLaunch(launch)

	// генерим хеш инстала
	if install.Id != 0 {

		if json_install, err := json.Marshal(install); err == nil {

			encoded := base64.StdEncoding.EncodeToString(json_install)
			response.InstallId = encoded + ":" + fmt.Sprint(launch.Launch_id)
		} else {
			fmt.Sprintln("install json marshaling error: ", err)
		}
	}

	return &response
}
