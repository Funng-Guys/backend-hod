package main

import (
    "fmt"
    "log"
    "mime"
    "net/http"
)

const authToken string = "2efhWdawJNO24THU9D2WQ"

func enforceJSONHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        contentType := r.Header.Get("Content-Type")

        if contentType == "" {
            http.Error(w, "Missing Content-Type header", http.StatusBadRequest)
            return
        }

        mt, _, err := mime.ParseMediaType(contentType)
		
        if err != nil {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
            return
		}

        if mt != "application/json" {
            http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func authHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusBadRequest)
            return
        }

        if authHeader != fmt.Sprintf("Bearer %s", authToken) {
            http.Error(w, "Invalid auth token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func provisionServer(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}

func main() {
    mux := http.NewServeMux()
    provisionHandler := http.HandlerFunc(provisionServer)
    mux.Handle("/", enforceJSONHandler(authHandler(provisionHandler)))

    log.Print("Listening on :3000...")
    err := http.ListenAndServe(":3000", mux)
    log.Fatal(err)
}		
