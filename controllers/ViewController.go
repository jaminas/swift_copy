package controllers

import (
	"github.com/astaxie/beego"
	"models"
	"service"
)

type ViewController struct {
	beego.Controller
	RequestService    *service.RequestService
	ClientService     *service.ClientService
	InstallService    *service.InstallService
	ClickhouseService *service.ClickhouseService
	MacrosService     *service.MacrosService
}

// @Title Get
// @Description find object by objectid
// @Param	Install	path	string	true	"Инсталл, полученный через метод /launch/..."
// @Success 302 Успешный редирект на партнера
// @Failure 400 Невалидные параметры запроса либо webview_allow=false
// @Failure 500 Ошибка сервиса
// @router /view/:install [get]
func (this *ViewController) Get() {

	response := models.ResponseView{}

	// парсим реквест на входящие данные и валидируем
	request := this.RequestService.ParseView(this.Ctx, this.GetString(":install"))
	response.Request = request

	// парсим данные клиента
	client := this.ClientService.Parse(this.Ctx)
	response.Client = client

	// ищем инсталл в кеше или базе
	install := this.InstallService.Get(request.InstallObject.Id)
	if install.Id == 0 { // инсталл может быть не найден в кеше или базе если это другая нода или не произошла синхронизация с базой

		// восстанавливаем инсталл из запроса - именно поэтому мы передаем объект инстала в запросе
		install = request.InstallObject
	}
	response.Install = install

	// генерим событие открытия webview
	event_view := models.NewEventView(install, request, client)
	response.EventView = event_view

	// добавляем событие открытия webview в КХ
	this.ClickhouseService.AddView(event_view)

	// облогораживаем макросами
	response.RedirectUrl = this.MacrosService.FillMacros(install.Webview_url, install, client, event_view)

	// в случае запрета вебвью возвращаем 400
	if install.Allow_webview == false {

		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	} else {

		this.Redirect(response.RedirectUrl, 302)
		return
	}

}
