package _28204

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("Hello, World!"))
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
