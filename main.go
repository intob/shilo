package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("set a PORT env var")
	}

	addr := ":" + port
	log.Println("listening on", addr)

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}
