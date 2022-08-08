package main

import (
	"fmt"
	"net/http"
	"testing"
)

// TestNoSurf - тестирует посредник NoSurf
func TestNoSurf(t *testing.T) {

	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// ничего не делаем если ошибки нет
	default:
		t.Error(
			fmt.Sprintf("Возращаемый формат не является http.Handler, его формат %T", v),
		)
	}

}

// TestNoSurf - тестирует посредник SessionLoad
func TestSessionLoad(t *testing.T) {

	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// ничего не делаем если ошибки нет
	default:
		t.Error(
			fmt.Sprintf("Возращаемый формат не является http.Handler, его формат %T", v),
		)
	}
}
