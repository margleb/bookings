package helpers

import (
	"fmt"
	"github.com/margleb/booking/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers - инициализирует новый хелпер
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError - клиентские ошибки
func ClientError(w http.ResponseWriter, status int) {
	// пишем ошибку в консколь
	app.InfoLog.Println("Client error with status of", status)
	// пишем ошибку в браузер
	http.Error(w, http.StatusText(status), status)
}

// ServerError - ceрверные ошибки
func ServerError(w http.ResponseWriter, err error) {
	// создаем тресировку кода
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// пишем тресировку серверной ошибки в консоль
	app.ErrorLog.Println(trace)
	// пишем ошибку в браузер
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
