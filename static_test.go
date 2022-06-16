package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"jgt.solutions/controllers"
)

func TestGETIndex(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	controllers.NewStatic().NewHome(response, request)

	got := response.Body.String()
	teamTitle := "Our Amazing Team"
	contactTitle := "Contact Us"

	if !strings.Contains(got, teamTitle) {
		t.Errorf("Can not find title %s", teamTitle)
	} else if !strings.Contains(got, contactTitle) {
		t.Errorf("Can not find title %s", contactTitle)
	}
}
