package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Henry!")
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
func envHandler(w http.ResponseWriter, r *http.Request) {
	for _, pair := range os.Environ() {
		fmt.Fprintln(w, pair)
	}
}
func argsHandler(w http.ResponseWriter, r *http.Request) {
	for i, arg := range os.Args {
		fmt.Fprintf(w, "[%02d] %s\n", i, arg)
	}
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/argz", argsHandler)
	http.HandleFunc("/varz", envHandler)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", getEnv("PORT", "8080")),
			nil,
		),
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
