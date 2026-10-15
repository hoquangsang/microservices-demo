package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://hello:hello@hellodb:5432/hello?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		var msg string
		err := db.QueryRowContext(r.Context(), "SELECT message FROM greetings ORDER BY id LIMIT 1").Scan(&msg)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			log.Printf("query: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "db error"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": msg})
	})
	http.HandleFunc("/_healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "db not ready", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ok"))
	})

	log.Printf("helloservice listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
