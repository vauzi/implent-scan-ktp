package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/vauzi/implent-scan-ktp/web"
)

const (
	port = ":9000"
)

func main() {

	mux := web.Route()
	fmt.Printf("Starting server at port on %s\n", port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) // Be careful with '*', it allows all origins
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Apply the CORS middleware
	handler := handlers.CORS(originsOk, headersOk, methodsOk)(mux)

	err := http.ListenAndServe(port, handler)

	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
