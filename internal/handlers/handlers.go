package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/driver"
	"github.com/margleb/booking/internal/forms"
	"github.com/margleb/booking/internal/helpers"
	"github.com/margleb/booking/internal/models"
	"github.com/margleb/booking/internal/render"
	"github.com/margleb/booking/internal/repository"
	"github.com/margleb/booking/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRep(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	// получаем Reservation значения из сесси
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	// получаем информацию о комнате
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// указываем имя комнаты которая забронирована
	res.Room.RoomName = room.RoomName

	// делаем кастинг обратно в строку
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil), // по умолчанию ошибок валидации нет
		Data:      data,           // при get запросе передаем пустое значение
		StringMap: stringMap,
	})
}

// PostReservation - пост запрос из формы
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// если не получается спарсить данные формы
	err := r.ParseForm()
	if err != nil {
		// запускаем сервеную ошибку
		helpers.ServerError(w, err)
		return
	}

	// конвертируем start/end date в правильный формат
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2006-01-02 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// конвертируем room_id в int формат
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Reservation - сохраняем данные из формы
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
	}

	// Возращаем новый поинтер формы
	form := forms.New(r.PostForm)

	// проверям, не пустое ли значение first_name
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone")
	// мин длина имени - три символа
	form.MinLength("first_name", 3)
	// проверяем, является ли email значение
	form.IsEmail("email")

	// если есть ошибки валидации
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form, // по умолчанию ошибок валидации нет
			Data: data,
		})
		return
	}

	// добавляем в базу данных
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomId,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	// добавляем ограничение по данной комнате
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// добавляем в сессию данные бронирования для отображения на странице результатов бронирования
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// делаем редирект на страницу /reservation-main
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// ReservationSummary - страница результатов бронирования
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error form session")
		// если не получилось взять данные из формы, то добавляем ошибку в сессию
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		// а также добавляем редирект на главную
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// удаляем из сессии reservation
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the about page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the about page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// конвертируем в time.Time
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// получаем список доступных комнат
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// выводим циклом доступные комнаты
	// for _, i := range rooms {
	//	m.App.InfoLog.Println("Room:", i.ID, i.RoomName)
	//}

	// если нет ни одной доступной комнаты
	if len(rooms) == 0 {
		// m.App.InfoLog.Println("No Avail")
		// Выводим сообщение о том, что нет свободных номеров
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// если есть свободные комнаты
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// сохраняем в сессию даты, чтобы передать на страницу бронирования
	m.App.Session.Put(r.Context(), "reservation", res)

	// передаем их в шаблон
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

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
		helpers.ServerError(w, err)
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
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	// получаем id комнаты из ссылки, конвертируем строку в int
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// m.App.Session.Get(r.Context(), "reservation")

	// получаем Reservation значения из сесси
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	// добавляем id комнаты в reservation
	res.RoomID = roomID

	// помещаем снова reservation в сессию
	m.App.Session.Put(r.Context(), "reservation", res)

	// делаем редирект на make-reservation, но уже с датами из сессии и id комнаты
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
