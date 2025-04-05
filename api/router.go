package api

import "github.com/go-chi/chi/v5"

func SetupRoutes(r chi.Router, handler *TaskHandler) {
	r.Route("/api/v1/task", func(r chi.Router) {
		r.Post("/", handler.CreateTask)
		r.Get("/", handler.ListTasks)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetTask)
			r.Put("/", handler.UpdateTask)
			r.Delete("/", handler.DeleteTask)
		})
	})
}
