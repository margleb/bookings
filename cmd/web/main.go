package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/driver"
	"github.com/margleb/booking/internal/handlers"
	"github.com/margleb/booking/internal/helpers"
	"github.com/margleb/booking/internal/models"
	"github.com/margleb/booking/internal/render"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// не может быть изменяемым
const portNumber = ":8080"

// настройки приложения
var app config.AppConfig

// сессии (указываются в конфиге)
var session *scs.SessionManager

// клиентские ошибки
var InfoLog *log.Logger

// серверные ошибки
var ErrorLog *log.Logger

// main is the main application function
func main() {

	// пробуем запустить приложение
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	// отложенное завершение подключения к БД
	log.Println("Connected to database")
	defer db.SQL.Close()

	// закрываем канал после отправки письма
	defer close(app.MailChan)

	// запускаем слушателя email

	fmt.Println("Starting mail listener")
	listenForMail()

	from := "me@here.com"
	auth := smtp.PlainAuth("", from, "", "localhost")
	err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello, world"))
	if err != nil {
		log.Println(err)
	}

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

// run - функция позволяющая проводить тестирование
func run() (*driver.DB, error) {
	// уточняем какого типа данные мы хотим хранить в сессии
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// канал созданный для отправки писем
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false // находится ли сайт в продакшене

	// инициализация клиентского лога
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = InfoLog

	// инициализация серверного лога
	ErrorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog

	session = scs.New()                            // новая сессия
	session.Lifetime = 24 * time.Hour              // 24 часа
	session.Cookie.Persist = true                  // должны ли cookie оставаться после закрытия
	session.Cookie.SameSite = http.SameSiteLaxMode // куки применяются к текущему сайту
	session.Cookie.Secure = app.InProduction       // ну ли криптовать куки

	app.Session = session // устанавливаем в конфиг сессии

	// устанавливаем соединение с базой данных
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=marglebdm")
	if err != nil {
		log.Fatal("Cannot connect to database. Dying...")
		return nil, err
	}

	// получаем кеш шаблонов
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create tmp cache")
		return nil, err
	}
	// сохраняем его в гл. переменную TemplateCache
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
