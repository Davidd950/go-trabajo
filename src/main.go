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

mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})


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
		w.Write([]byte(html))
	})

	log.Printf("Servidor iniciado en puerto %s", port)
	handler := loggingMiddleware(mux)
	log.Fatal(http.ListenAndServe(":"+port, handler))

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

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ww := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(ww, r)
		_ = time.Since(start)

		logRequest(r, ww.statusCode)
	})
}
