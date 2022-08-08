package main

import "testing"

// TestRun - позволяет тестировать осн. функцию приложения
func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Error("failed run()")
	}
}
