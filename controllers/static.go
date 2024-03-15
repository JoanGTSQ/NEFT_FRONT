package controllers

import (

	"jgt.solutions/views"
)

func NewStatic() *Static {
	return &Static{
		NotFound: views.NewView("dashboard", "static/404"),
		Error:    views.NewView("dashboard", "static/505"),
	}
}

type Static struct {
	NotFound *views.View
	Error    *views.View
}

