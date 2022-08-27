package render

import (
	"github.com/margleb/booking/internal/models"
	"net/http"
	"testing"
)

// TestAddDefaultData - тестирует добавление данных по умолчанию
func TestAddDefaultData(t *testing.T) {

	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// пробуем положить флеш в тек
	session.Put(r.Context(), "flash", "123")

	// запускаем функцию тестирования
	result := AddDefaultData(&td, r)

	// если не найден в templateData 123, то выводим ошибку
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

// getSession - создает реквест с сессией
func getSession() (*http.Request, error) {

	// создаем реквест
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// создает сессию с текущим контекстом
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-session"))

	// кладем в реквест текущий контекст с сессией
	r = r.WithContext(ctx)

	return r, nil
}

func TestRenderTemplate(t *testing.T) {

	// обязательно переписываем пути относительно папки
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	// реквест для RenderTemplate
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// ResponseWriter c необходимыми интерфесами
	var ww myWriter

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "not-existed.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("has to be error but haven't")
	}

}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {

	// обязательно переписываем пути относительно папки
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error("creating tmpl cache failed")
	}
}
