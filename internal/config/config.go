package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/margleb/booking/internal/models"
	"log"
	"text/template"
)

// конфиг настройки приложения
type AppConfig struct {
	UseCache      bool                          // нужно ли использовать кеш
	TemplateCache map[string]*template.Template // кеш шаблонов
	InfoLog       *log.Logger                   // клиентские ошибки
	ErrorLog      *log.Logger                   // серверные ошибки
	InProduction  bool                          // находится ли сайт в продакшене
	Session       *scs.SessionManager
	MailChan      chan models.MailData // канал для писем
}
