package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

func NewRouter(handler *handler.ApiHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// For Desktop
	r.Post("employee/auth", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("screenshot/{employeeId}", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("statistic/{employeeId}/", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("work_session/{employeeId}/start", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("work_session/{employeeId}/stop", func(w http.ResponseWriter, r *http.Request) {})

	// For Admin Panel
	r.Post("manager/auth", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("screenshot/{employeeId}/{date}", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("statistic/{employeeId}/{date}", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("work_session/{employeeId}/{date}", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("employee", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("photo", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("manager/{managerId}/employee/all", func(w http.ResponseWriter, r *http.Request) {})
	r.Delete("employee/{employeeId}", func(w http.ResponseWriter, r *http.Request) {})
	return r
}
