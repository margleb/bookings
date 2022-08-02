package forms

// тип данных ошибки
type errors map[string][]string

// Add - добавляет ошибки для поля в slice
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get - получает сообщение об ошибке
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	} else {
		return es[0]
	}
}
