package main

import (
	"evilllm-http-api/routes"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("x-real-ip")
		if ip == "" {
			ip = r.Header.Get("x-forwarded-for")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		log.Info("Request received", "to url", r.URL, "with method", r.Method, "from addr", ip)

		next.ServeHTTP(w, r)

		log.Info("Response written", "to url", r.URL, "with method", r.Method, "from addr", ip)
	})
}

func main() {
	log.SetLevel(log.DebugLevel)
	logger := log.NewWithOptions(os.Stderr, log.Options{Prefix: "http"})
	stdlog := logger.StandardLog(log.StandardLogOptions{
		ForceLevel: log.DebugLevel,
	})

	if err := os.MkdirAll("_UPLOADS_", os.ModePerm); err != nil {
		log.Error("Could not create uploads directory", "err", err)
	}
	if err := os.MkdirAll("_GENERATED_", os.ModePerm); err != nil {
		log.Error("Could not create uploads directory", "err", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", routes.HealthHandler)
	mux.HandleFunc("/health", routes.HealthHandler)
	mux.HandleFunc("/ai/communicate", routes.CommunicateHandler)

	loggedMux := requestLogger(mux)

	server := &http.Server{
		Addr: ":8080", ErrorLog: stdlog,
		Handler: loggedMux,
	}

	log.Info("Starting server on :8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Error("Could not start server: %v", "err", err)
	}
}
