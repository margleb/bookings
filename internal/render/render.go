package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/models"
	"net/http"
	"path/filepath"
	"text/template"
)

// массив функций применяемых к шаблону
var functions = template.FuncMap{}

var app *config.AppConfig

// статически прописанный путь до шаблонов
var pathToTemplates = "./templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData передает в шаблоны данные по умолчанию
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// по умолчанию генерируется CSRFToken
	td.CSRFToken = nosurf.Token(r)
	// а также значения для сессии
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	//  если необходимо использовать кеш
	if app.UseCache {
		// получаем значение кеша из ранее сгенерированного
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// берем шаблон из сгенерированного кеша
	t, ok := tc[tmpl]
	// если шаблон не найден
	if !ok {
		return errors.New("can't find cache template")
	}

	buf := new(bytes.Buffer)

	// в случаем если данные по умолчанию то добавляет их
	td = AddDefaultData(td, r)

	// записываем в буфер полученный шаблон
	_ = t.Execute(buf, td)
	// записываем его в ответ
	_, err := buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
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
