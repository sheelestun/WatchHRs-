package router

import (
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
	r.Post("/employee/auth", handler.AuthEmployeeHandler)
	r.Post("/screenshot/{employeeId}", handler.AddScreenshotHandler)
	r.Post("/statistic/{employeeId}", handler.AddScreenshotStatisticHandler)
	r.Post("/work_session/{employeeId}/start", handler.StartWorkSessionHandler)
	r.Post("/work_session/{employeeId}/stop", handler.StopWorkSessionHandler)

	// For Admin Panel
	r.Post("/manager/auth", handler.AuthManagerHandler)
	r.Get("/screenshot/{employeeId}/{date}", handler.GetScreenshotsHandler)
	r.Get("/statistic/{employeeId}/{date}", handler.GetScreenshotsStatisticHandler)
	r.Get("/work_session/{employeeId}/{date}", handler.GetWorkSessionsHandler)
	r.Post("/employee", handler.AddEmployeeInfoHandler)
	r.Post("/photo", handler.AddEmployeePhoto)
	r.Get("/manager/{managerId}/employee/all", handler.GetAllEmployeesInfoByManagerIDHandler)
	r.Delete("/employee/{employeeId}", handler.DeleteEmployee)
	return r
}
