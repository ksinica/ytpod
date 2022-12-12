package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ksinica/ytpod/pkg/youtube"
)

type transport struct{}

func (*transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set(
		"User-Agent",
		"github.com/ksinica/ytpod",
	)
	return http.DefaultTransport.RoundTrip(r)
}

func main() {
	c := http.Client{
		Transport: new(transport),
	}

	r := chi.NewMux()
	r.Use(youtube.UseHTTPClient(&c))
	r.Use(middleware.Timeout(time.Second * 60))

	r.Route("/youtube", youtube.Router)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/ksinica/ytpod", http.StatusSeeOther)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
