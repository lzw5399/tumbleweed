package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "bank/workflow/engine/src/docs"
	_ "bank/workflow/engine/src/initialize"
	"bank/workflow/engine/src/router"
	_ "bank/workflow/engine/src/router"
)

func main() {
	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("application has started, listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
