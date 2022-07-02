package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

// WriteToConsole - пишет информацию в консоль
func WriteToConsole(next http.Handler) http.Handler {
	// возращает обработчик
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// пишет в консоль
		fmt.Println("Hit the page")
		// продолжает выполнение скрипта
		next.ServeHTTP(w, r)
	})
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
