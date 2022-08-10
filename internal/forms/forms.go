package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
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

func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field) // получаем значение поля
		if strings.TrimSpace(value) == "" {
			// если значение пустое - добавляем ошибку
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has - проверят что поле не является пустым
func (f *Form) Has(field string) bool {
	x := f.Get(field) // получаем поле
	if x == "" {
		// добавляем сообщение об ошибке
		f.Errors.Add(field, "This field cannot bе blank")
		return false
	}
	return true
}

// MinLength - минимальная длина значений поля
func (f *Form) MinLength(field string, length int) bool {

	// 1. получаем значение из формы
	x := f.Get(field)
	// 2. если полученное значение меньше длины
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail - проверяет, является ли строка email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
