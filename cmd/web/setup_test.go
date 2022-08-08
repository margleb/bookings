package main

import (
	"net/http"
	"os"
	"testing"
)

// TestMain - запускается до тестирования, создаем http.Handler для middleware_test.go
func TestMain(m *testing.M) {
	// перед тем как запустить тест, что то выполняется в функции, а затем выходит
	os.Exit(m.Run())
}

// обьект необходимый для TestNoSurf и TestSessionLoad
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
