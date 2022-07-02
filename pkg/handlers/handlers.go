package handlers

import (
	"github.com/margleb/booking/pkg/config"
	"github.com/margleb/booking/pkg/models"
	"github.com/margleb/booking/pkg/render"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// берем ip адрес поль-ля
	remoteIP := r.RemoteAddr
	// помещаем его в config c сессиями
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// подготовка данных для передачи в шаблон
	stringMap := make(map[string]string)
	stringMap["example"] = "Hello, world"

	// получаем из конфига значение сессии
	// делаем явное преобразование к строке
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	// добавляем значение сессии в stringMap
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
