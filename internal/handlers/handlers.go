package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/forms"
	"github.com/margleb/booking/internal/models"
	"github.com/margleb/booking/internal/render"
	"log"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
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

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var reservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // по умолчанию ошибок валидации нет
		Data: data,           // при get запросе передаем пустое значение
	})
}

// PostReservation - пост запрос из формы
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// если не получается спарсить данные формы
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// Reservation - сохраняем данные из формы
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// Возращаем новый поинтер формы
	form := forms.New(r.PostForm)

	// проверям, не пустое ли значение first_name
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone")
	// мин длина имени - три символа
	form.MinLength("first_name", 3, r)

	// если есть ошибки валидации
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form, // по умолчанию ошибок валидации нет
			Data: data,
		})
		return
	}

}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the about page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the about page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	_, _ = w.Write([]byte(fmt.Sprintf("Started date: %s, ended date: %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON обрабатывает JSON запросы
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	// Пример JSON
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	// Преобразуем struct в JSON
	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return
	}

	// Выводим лог, а также устанавливаем заголовок
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")

	// Записываем в ответ
	_, _ = w.Write(out)

}

// Contact is the about page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
