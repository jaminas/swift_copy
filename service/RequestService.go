package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/speps/go-hashids"
	"models"
	"regexp"
	"strconv"
	"strings"
)

type RequestService struct {
	logger           *Logger
	hash_id          *hashids.HashID
	query_alphabet   *regexp.Regexp
	query_param_size int
}

/**
 * Конструктор
 */
func NewRequestService(logger *Logger, hash_id *hashids.HashID) *RequestService {
	return &RequestService{
		logger:           logger,
		hash_id:          hash_id,
		query_alphabet:   regexp.MustCompile("[^0-9A-Za-z_.,:;{}()\\[\\]\\-=+\\*|]+"),
		query_param_size: 64,
	}
}

/**
 * Метод парсит входящие данные и возвращает модель входящих данных
 */
func (this *RequestService) Parse(ctx *context.Context, tail string) *models.RequestLaunch {

	request := models.RequestLaunch{}

	if ctx.Input.RequestBody != nil {

		if decoded, err := base64.StdEncoding.DecodeString(string(ctx.Input.RequestBody)); err != nil {
			this.logger.Warn(fmt.Sprintf("Post body base64 decode error: %v", err))
		} else {

			if len(decoded) > 0 {
				if err := json.Unmarshal(decoded, &request); err != nil {
					this.logger.Warn(fmt.Sprintf("Post body Unmarshal json error: %v", err))
				}
			}
		}
	}

	request.CampaignId, request.DomainId = this.parseTail(tail)
	request.InstallObject, _ = this.parseInstall(request.Install)

	return &request
}

/**
 */
func (this *RequestService) ParseView(ctx *context.Context, tail string) *models.RequestView {

	request := models.RequestView{}
	request.InstallObject, request.LaunchId = this.parseInstall(tail)
	return &request
}

/**
 * Метод парсит хвост на кампанию и домен
 */
func (this *RequestService) parseTail(tail string) (int64, int64) {

	var (
		campaign_id int64
		domain_id   int64
	)

	if numbers, err := this.hash_id.DecodeInt64WithError(tail); err == nil {

		if len(numbers) == 4 {

			if numbers[0] == 1 { // первое число знак, т.к. отрицательные числа hashId не принимает
				campaign_id = numbers[1]
			} else {
				campaign_id = numbers[1] - numbers[1]*2
			}

			if numbers[2] == 1 { // первое число знак, т.к. отрицательные числа hashId не принимает
				domain_id = numbers[3]
			} else {
				domain_id = numbers[3] - numbers[3]*2
			}
		} else {
			this.logger.Warn("Tail encode error: hash consists more then 4 digits")
		}
	} else {
		this.logger.Warn(fmt.Sprintf("Launch hash decode error: %v", err))
	}

	return campaign_id, domain_id
}

/**
 * Метод парсит хэш install
 */
func (this *RequestService) parseInstall(param string) (*models.Install, int64) {

	if param != "" {

		part := strings.Split(param, ":")
		if len(part) == 2 {

			if decoded, err := base64.StdEncoding.DecodeString(part[0]); err == nil {

				if len(decoded) > 0 {

					install := models.Install{}
					if err := json.Unmarshal(decoded, &install); err == nil {

						var launch_id int64
						if n, err := strconv.Atoi(part[1]); err == nil {
							launch_id = int64(n)
						} else {
							this.logger.Warn(fmt.Sprintf("Install base64 decode error: %v", err))
						}

						return &install, launch_id

					} else {

						this.logger.Warn(fmt.Sprintf("Install unmarshal json error: %v", err))
					}
				}

			} else {
				this.logger.Warn(fmt.Sprintf("Install base64 decode error: %v", err))
			}
		} else {
			this.logger.Warn("Install parts must be 2")
		}
	}

	return nil, 0
}

/**
 * Метод нормализует строку согласно длительности и алфавита
 */
func (this *RequestService) normalizeString(param string) string {

	result := ""
	if strings.TrimSpace(param) != "" {

		result = this.query_alphabet.ReplaceAllString(param, "")
		if len([]rune(result)) > this.query_param_size {
			result = string([]rune(result)[:this.query_param_size])
		}
	}

	return result
}
