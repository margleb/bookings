package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/models"
	"github.com/margleb/booking/internal/render"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"text/template"
	"time"
)

// осн. настойки сайта
var app config.AppConfig

// сессии
var session *scs.SessionManager

// статический путь
var pathToTemplates = "./../../templates"

// функции передаваемые в шаблон
var functions = template.FuncMap{}

// TestMain - основная фун-ция запускаемая для тестов
func TestMain(m *testing.M) {

	// уточняем какого типа данные мы хотим хранить в сессии
	gob.Register(models.Reservation{})

	app.InProduction = false // находится ли сайт в продакшене

	// инициализация клиентского лога
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = InfoLog

	// инициализация серверного лога
	ErrorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog

	session = scs.New()                            // новая сессия
	session.Lifetime = 24 * time.Hour              // 24 часа
	session.Cookie.Persist = true                  // должны ли cookie оставаться после закрытия
	session.Cookie.SameSite = http.SameSiteLaxMode // куки применяются к текущему сайту
	session.Cookie.Secure = app.InProduction       // ну ли криптовать куки

	app.Session = session // устанавливаем в конфиг сессии

	// получаем кеш шаблонов
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create tmp cache")

	}
	// Сохраняем его в гл. переменную TemplateCache
	app.TemplateCache = tc
	app.UseCache = false // не используем кеш данных

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())

}

// getRoutes - возращает маршруты
func getRoutes() http.Handler {

	// маршрут
	mux := chi.NewRouter()
	// посредник - поглащает панику и печает стек
	mux.Use(middleware.Recoverer)
	// посредник - просто пишет в консоль
	// mux.Use(WriteToConsole)
	// посредник - устанавливаем CSRF токен
	// mux.Use(NoSurf) - СSRF токены в тестах не нужны
	// посредник - позволяет использовать сессии
	mux.Use(SessionLoad)

	// уст. маршруты
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	// создаем сервер файлов
	fileServer := http.FileServer(http.Dir("./static/"))
	// помещаем ее в Handle
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// возвращаем маршруты
	return mux

}

// NoSurf - позволяет работать с CSRF токенами
func NoSurf(next http.Handler) http.Handler {
	// обработчик csrf
	csrfHandler := nosurf.New(next)
	// устанавливаем в обрабатчики куки
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction, // защиенные куки
		SameSite: http.SameSiteLaxMode,
	})
	// возращаем обработчик
	return csrfHandler
}

// SessionLoad - указываем, что необходимо использовать сессии
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTemplateCache создаем кеш шаблона как map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// создаем map кеш страниц
	myCache := map[string]*template.Template{}

	// сканируем все страницы *.page.tmp
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// перебираем кажду страницу
	for _, page := range pages {

		// название страницы
		name := filepath.Base(page)

		// парсим конкретную страницу
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// сканируем все страницы c *.layouts.tmp
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// если есть шаблоны, то парсим их
		if len(matches) > 0 {
			// парсим к странице конкретный шаблон
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		// добавляем их в кеш
		myCache[name] = ts

	}

	return myCache, nil

}
