// main.go
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
        <html>
        <head><title>Despliegue Seguro</title></head>
        <body>
            <h1>Sistema de Despliegue Seguro</h1>
            <img src="/static/logo.png" alt="Logo" />
        </body>
        </html>`
		w.Write([]byte(html))
	})

	log.Printf("Servidor iniciado en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
