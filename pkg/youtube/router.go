package youtube

import "github.com/go-chi/chi/v5"

func Router(r chi.Router) {
	r.Get("/feed/*", handleGetFeed)
	r.Get("/stream/{videoID}", handleGetStream)
}
