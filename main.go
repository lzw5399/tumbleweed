package main

import (
	"bank/distributedquery/src/router"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "bank/distributedquery/src/docs"
	_ "bank/distributedquery/src/initialize"
	_ "bank/distributedquery/src/router"
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
