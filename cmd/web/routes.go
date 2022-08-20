package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/handlers"
	"net/http"
)

// Routes - маршруты с мультиплексером
func routes(app *config.AppConfig) http.Handler {
	// маршрут
	mux := chi.NewRouter()
	// посредник - поглащает панику и печает стек
	mux.Use(middleware.Recoverer)
	// посредник - просто пишет в консоль
	// mux.Use(WriteToConsole)
	// посредник - устанавливаем CSRF токен
	mux.Use(NoSurf)
	// посредник - позволяет использовать сессии
	mux.Use(SessionLoad)

	// уст. маршруты
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	// создаем сервер файлов
	fileServer := http.FileServer(http.Dir("./static/"))
	// помещаем ее в Handle
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// возращаем
	return mux
}
