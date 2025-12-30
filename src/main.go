package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
	<html>
	<head><title>Despliegue Seguro</title></head>
	<body>
		<h1>Sistema de Despliegue Seguro</h1>
		<img src="/static/logo.png" alt="Logo" />
		<p>Este es un texto de prueba</p>
	</body>
	</html>`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})

	log.Printf("Servidor iniciado en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, loggingMiddleware(mux)))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Error interno simulado", http.StatusInternalServerError)
	})

}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status,omitempty"`
}

func logRequest(r *http.Request, status int) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "INFO",
		Message:   "HTTP Request",
		Path:      r.URL.Path,
		Status:    status,
	}
	jsonLog, _ := json.Marshal(entry)
	log.Println(string(jsonLog))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rr, r)

		logRequest(r, rr.status)

		log.Printf(
			`{"level":"INFO","message":"request processed","method":"%s","path":"%s","duration_ms":%d}`,
			r.Method,
			r.URL.Path,
			time.Since(start).Milliseconds(),
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}
