package models

import "github.com/margleb/booking/internal/forms"

// TemplateData Данные для передачи в шаблон
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]string
	FloatMap  map[string]float32
	Data      map[string]interface{} // передача неизвестного типа данных
	CSRFToken string                 // токет межсайтовой подделки для форм
	Flash     string                 // флеш сообщения
	Warning   string                 // предупреждения
	Error     string                 // ошибки
	Form      *forms.Form
}
