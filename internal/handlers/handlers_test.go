package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// данные отправленные постом (ключ, значение)
type postData struct {
	key   string
	value string
}

// тестируемые маршруты
var theTests = []struct {
	name               string     // название маршрута
	url                string     // путь до маршрута
	method             string     // используемый метод в маршруте
	params             []postData // передаваемые параметры
	expectedStatusCode int        // ожидаемый статус ответа
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"ma", "/make-reservation", "GET", []postData{}, http.StatusOK},
}

// тестирование обработчиков
func TestHandlers(t *testing.T) {
	// 1. получаем маршруты
	routes := getRoutes()

	// 2. создаем тестовый сервер
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() // закрываем его

	// 3. перебираем каждый тест
	for _, e := range theTests {
		// 3.1 для поста и гет запросов
		if e.method == "GET" {
			// получаем результат гет запроса
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			// если полученный статус кода не соответсвуте ожидаемому
			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}

		} else {

		}
	}

}
