package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

// сессии
var session *scs.SessionManager

// тестовая переменная настроек
var testApp config.AppConfig

// TestMain - вызывается перед другими функциями
func TestMain(m *testing.M) {

	// уточняем какого типа данные мы хотим хранить в сессии
	gob.Register(models.Reservation{})

	testApp.InProduction = false // находится ли сайт в продакшене

	session = scs.New()                            // новая сессия
	session.Lifetime = 24 * time.Hour              // 24 часа
	session.Cookie.Persist = true                  // должны ли cookie оставаться после закрытия
	session.Cookie.SameSite = http.SameSiteLaxMode // куки применяются к текущему сайту
	session.Cookie.Secure = false                  // ну ли криптовать куки

	testApp.Session = session // устанавливаем в конфиг сессии

	// ссылка на тестовую переменную настроек в саму render.go
	app = &testApp

	// запускаем тесты, перед закрытием
	os.Exit(m.Run())
}
