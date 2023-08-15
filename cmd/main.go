package main

import (
	"log"
	"net/http"
	"os"

	"server/pkg/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	//make new router
	router := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	//router list
	router.HandleFunc("/rekomendasi", controllers.GetRekomendasi).Methods(http.MethodPost)
	router.HandleFunc("/", controllers.Channels).Methods(http.MethodGet)

	log.Println("Api is in 4000")
	// controllers.Kelompokin()

	http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(headers, methods, origins)(router))
}
