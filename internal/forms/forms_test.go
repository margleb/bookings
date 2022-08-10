package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {

	// создаем тестовый реквест
	r := httptest.NewRequest("POST", "/test-url", nil)

	// создаем поинтер формы
	form := New(r.PostForm)

	// пробуем создать валидацию для несуществующих полей
	// должен добавить ошибки в Form.Errors
	form.Required("a", "b", "с")

	// не должна быть валидна, так как необходимых полей нет
	if form.Valid() {
		t.Error("form shows valid when required is missing")
	}

}

func TestForm_Required(t *testing.T) {

	// создаем тестовый реквест
	r := httptest.NewRequest("POST", "/test-url", nil)

	// создаем поинтер формы
	form := New(r.PostForm)

	// пробуем создать валидацию для несуществующих полей
	// должен добавить ошибки в Form.Errors
	form.Required("a", "b", "с")

	// не должна быть валидна, так как необходимых полей нет
	if form.Valid() {
		t.Error("form shows valid when required is missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	// создаем тестовый реквест
	r = httptest.NewRequest("POST", "/test-url", nil)
	// добавляем данные в реквест
	r.PostForm = postedData
	// добавляем в объект данные из формы
	form = New(r.PostForm)

	// пробуем создать валидацию для несуществующих полей
	// должен добавить ошибки в Form.Errors
	form.Required("a", "b", "c")

	// не должна быть валидна, так как необходимых полей нет
	if !form.Valid() {
		t.Error("form shows invalid when required fields is exits")
	}

}

func TestForm_Has(t *testing.T) {
	// создаем реквест, в который будет передаваться данные
	r := httptest.NewRequest("POST", "/some-url", nil)
	// создаем новый поинтер на форму
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	form = New(postData)

	has = form.Has("a")
	if !has {
		t.Error("show form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	// создаем реквест, в который будет передаваться данные
	r := httptest.NewRequest("POST", "/some-url", nil)
	// создаем новый поинтер на форму
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// создаем значения
	postValues := url.Values{}
	postValues.Add("some_field", "some value")
	form = New(postValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("show minlength of 100 met when data is shorter")
	}

	postValues = url.Values{}
	postValues.Add("another_field", "abc123")
	form = New(postValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows min length of 1 is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	// создаем реквест, в который будет передаваться данные
	postValues := url.Values{}
	// создаем новый поинтер на форму
	form := New(postValues)

	form.IsEmail("x")

	if form.Valid() {
		t.Error("form shows valid email of non-existed field")
	}

	postValues = url.Values{}
	postValues.Add("email", "margleb.dm@yandex.ru")
	form = New(postValues)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("Got an invalid email when we should not have")
	}

	postValues = url.Values{}
	postValues.Add("email", "x")
	form = New(postValues)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("Got valid for invalid email address")
	}

}
