package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/A-Victory/blog-API/routes"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	r := httprouter.New()
	routes.Routers(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	fmt.Println("Starting server on port: " + port)
	log.Fatal(s.ListenAndServe())
}
