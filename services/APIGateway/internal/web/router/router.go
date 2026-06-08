package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

func NewRouter(authHandler *handler.AuthHandler, employeeHandler *handler.EmployeeHandler,
	imageHandler *handler.ImageHandler, managerHandler *handler.ManagerHandler,
	statisticHandler *handler.ScreenshotStatisticHandler, sessionHandler *handler.WorkSessionHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- HealthCheck ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// --- Public ---
	r.Post("/auth", authHandler.AuthHandler)
	r.Post("/manager", managerHandler.AddManagerInfoHandler)
	r.Post("/refresh", authHandler.RefreshTokenHandler)

	// --- Protected ---
	r.Group(func(r chi.Router) {
		r.Use(authHandler.JWTMiddleware)
		r.Use(handler.RequireRole("employee"))

		// Desktop
		r.Post("/screenshot/{employeeId}", imageHandler.AddScreenshotHandler)
		r.Post("/statistic/{employeeId}", statisticHandler.AddScreenshotStatisticHandler)
		r.Post("/work_session/{employeeId}/start", sessionHandler.StartWorkSessionHandler)
		r.Post("/work_session/{employeeId}/stop", sessionHandler.StopWorkSessionHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(authHandler.JWTMiddleware)
		r.Use(handler.RequireRole("manager"))

		// Admin
		r.Get("/screenshot/{employeeId}/{date}/archive", imageHandler.GetScreenshotsArchiveHandler)
		r.Get("/screenshot/{employeeId}/file/{filename}", imageHandler.GetScreenshotFileHandler)
		r.Get("/screenshot/{employeeId}/{date}", imageHandler.GetScreenshotsHandler)
		r.Get("/statistic/{employeeId}/{date}", statisticHandler.GetScreenshotsStatisticHandler)
		r.Get("/work_session/{employeeId}/{date}", sessionHandler.GetWorkSessionsHandler)
		r.Post("/employee", employeeHandler.AddEmployeeInfoHandler)
		r.Post("/photo", employeeHandler.AddEmployeePhoto)
		r.Get("/manager/{managerId}/employee/all", employeeHandler.GetAllEmployeesInfoByManagerIDHandler)
		r.Delete("/employee/{employeeId}", employeeHandler.DeleteEmployee)
	})
	return r
}
