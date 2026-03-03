package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

func NewRouter(apiHandler *handler.ApiHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- Public ---
	r.Post("/auth", apiHandler.AuthHandler)
	r.Post("/manager", apiHandler.AddManagerInfoHandler)
	r.Post("/refresh", apiHandler.RefreshTokenHandler)

	// --- Protected ---
	r.Group(func(r chi.Router) {
		r.Use(apiHandler.JWTMiddleware)
		r.Use(handler.RequireRole("employee"))

		// Desktop
		r.Post("/screenshot/{employeeId}", apiHandler.AddScreenshotHandler)
		r.Post("/statistic/{employeeId}", apiHandler.AddScreenshotStatisticHandler)
		r.Post("/work_session/{employeeId}/start", apiHandler.StartWorkSessionHandler)
		r.Post("/work_session/{employeeId}/stop", apiHandler.StopWorkSessionHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(apiHandler.JWTMiddleware)
		r.Use(handler.RequireRole("manager"))

		// Admin
		r.Get("/screenshot/{employeeId}/{date}", apiHandler.GetScreenshotsHandler)
		r.Get("/statistic/{employeeId}/{date}", apiHandler.GetScreenshotsStatisticHandler)
		r.Get("/work_session/{employeeId}/{date}", apiHandler.GetWorkSessionsHandler)
		r.Post("/employee", apiHandler.AddEmployeeInfoHandler)
		r.Post("/photo", apiHandler.AddEmployeePhoto)
		r.Get("/manager/{managerId}/employee/all", apiHandler.GetAllEmployeesInfoByManagerIDHandler)
		r.Delete("/employee/{employeeId}", apiHandler.DeleteEmployee)
	})
	return r
}
