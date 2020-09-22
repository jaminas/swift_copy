package controllers

import (
	"github.com/astaxie/beego"
	"models"
	"service"
)

type LaunchController struct {
	beego.Controller
	RequestService *service.RequestService
	ClientService  *service.ClientService
	BusinessLayer  *service.BusinessLayer
}

// @Title Get
// @Description Метод возвращает или создает инсталл
// @Param	tail	path	string  true	"Хэш кампании и домена, полученный приложением из кампании или deeplink"
// @Param	body	body	string  false	"Base64 json объект данных"
// @Success 200 {object} models.ResponseLaunch
// @Failure 400 Невалидные параметры запроса
// @Failure 500 Ошибка сервиса
// @router /launch/:tail [post]
func (this *LaunchController) Post() {

	// парсим реквест на входящие данные и валидируем
	request := this.RequestService.Parse(this.Ctx, this.GetString(":tail"))
	// определяем клиента
	client := this.ClientService.Parse(this.Ctx)

	business := this.BusinessLayer.Get(request, client)

	response := models.ResponseLaunch{}

	if business.Install != nil {
		response.InstallId = business.InstallId
		response.AllowWebview = business.Install.Allow_webview
	}

	this.Data["json"] = response
	this.ServeJSON()
}
