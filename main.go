package main

import (
	"log"
	"net/http"
	"os"

	// _ "bank/workflow/engine/src/docs"
	_ "workflow/src/initialize"
	"workflow/src/router"
)

func main() {
	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("监听的端口为: %s", port)
	log.Fatalf("应用启动失败，原因: %s\n", http.ListenAndServe(":"+port, r).Error())
}
