package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/margleb/booking/internal/config"
	"testing"
)

// TestRoutes - позволяет тестировать маршруты
func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// ничего не делаем
	default:
		t.Error(
			fmt.Sprintf("Возращаемый формат не является http.Handler, его формат %T", v),
		)
	}

}
