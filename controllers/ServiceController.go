package controllers

import (
	"github.com/astaxie/beego"
	"service"
)

type ServiceController struct {
	beego.Controller
	CampaignService *service.CampaignService
	PayService      *service.PayService
}

func (this *ServiceController) Get() {

	this.CampaignService.Load2Store()
	this.PayService.Load2Store()

	this.Data["json"] = "ok"
	this.ServeJSON()
}
