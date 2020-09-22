package service

import (
	"models"
	_ "time"
)

type ModuleService struct {
	logger *Logger
}

func NewModuleService(logger *Logger) *ModuleService {
	return &ModuleService{logger: logger}
}

/**
 * Метод валидирует запрос по модулям
 */
func (this *ModuleService) Check(request *models.RequestLaunch, client *models.Client, campaign *models.Campaign, stream *models.CampaignStream) bool {

	return true
}
