package config

import (
	"github.com/alexedwards/scs/v2"
	"text/template"
)

// конфиг настройки приложения
type AppConfig struct {
	UseCache      bool                          // нужно ли использовать кеш
	TemplateCache map[string]*template.Template // кеш шаблонов
	InProduction  bool                          // находится ли сайт в продакшене
	Session       *scs.SessionManager
}
