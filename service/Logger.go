package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type Logger struct {
	BeeLogger *logs.BeeLogger
}

/**
 * Конструктор
 */
func NewLogger() *Logger {

	BeeLogger := logs.NewLogger(10000)
	BeeLogger.SetLogger(logs.AdapterConsole, `{"level":1}`)
	BeeLogger.Async()

	return &Logger{BeeLogger: BeeLogger}
}

func (this *Logger) Warn(message string) {

	fmt.Println(message)
	this.slack(message)
	this.BeeLogger.Warn(message)
}

func (this *Logger) Info(message string) {

	fmt.Println(message)
	this.BeeLogger.Info(message)
}

func (this *Logger) slack(msg string) {

	request_body, err := json.Marshal(map[string]string{
		"text": "test_message",
	})
	if err != nil {
		fmt.Println("Marshaling logger message error: %s", err)
	}

	resp, err := http.Post(beego.AppConfig.String("slack_webhook"), "application/json", bytes.NewBuffer(request_body))
	if err != nil {
		fmt.Println("Sending message to slack error: %s", err)
	}
	defer resp.Body.Close()

}
