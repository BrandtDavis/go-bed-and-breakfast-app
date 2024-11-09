package main

import (
	"bookings/internal/config"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// do nothing; test passed
	default:
		t.Errorf("type is not *chi.Mux, instead it is %T", v)
	}
}
