package main

import (
	"log"
	"net/http"
	"os"

	// _ "bank/workflow/engine/src/docs"
	_ "workflow/initialize"
	"workflow/router"
	_ "workflow/router"
)

func main() {
	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("application has started, listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
