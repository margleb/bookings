package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/handlers"
	"github.com/margleb/booking/internal/models"
	"github.com/margleb/booking/internal/render"
	"log"
	"net/http"
	"time"
)

// не может быть изменяемым
const portNumber = ":8080"

// настройки приложения
var app config.AppConfig

// сессии (указываются в конфиге)
var session *scs.SessionManager

// main is the main application function
func main() {

	// уточняем какого типа данные мы хотим хранить в сессии
	gob.Register(models.Reservation{})

	app.InProduction = false // находится ли сайт в продакшене

	session = scs.New()                            // новая сессия
	session.Lifetime = 24 * time.Hour              // 24 часа
	session.Cookie.Persist = true                  // должны ли cookie оставаться после закрытия
	session.Cookie.SameSite = http.SameSiteLaxMode // куки применяются к текущему сайту
	session.Cookie.Secure = app.InProduction       // ну ли криптовать куки

	app.Session = session // устанавливаем в конфиг сессии

	// получаем кеш шаблонов
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create tmp cache")

	}
	// сохраняем его в гл. переменную TemplateCache
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// создаем сервер
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// запускаем его
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
