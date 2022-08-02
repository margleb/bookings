package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
	url.Values        // добавляет значения полей
	Errors     errors // ошибки
}

// New - фун-ция инициализации новой формы
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has - проверят что поле не является пустым
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field) // получаем поле
	if x == "" {
		return false
	}
	return true
}
