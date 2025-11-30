package tests_test

import "testing"

func TestFailing(t *testing.T) {
	t.Error("failed succesfully")
}
