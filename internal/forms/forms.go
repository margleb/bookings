package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
	url.Values        // добавляет значения полей
	Errors     errors // ошибки
}

// Valid возращает true, если нет ошибок, иначе false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New - фун-ция инициализации новой пустой формы формы
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
		// добавляем сообщение об ошибке
		f.Errors.Add(field, "This field cannot bе blank")
		return false
	}
	return true
}
