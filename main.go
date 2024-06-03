package main

import (
	"entale-test/models"
	"fmt"
	"net/http"

	"entale-test/controllers"

	"github.com/go-chi/chi/v5"
)


func main() {
	models.ConnectDataBase()
	if models.DB == nil {
		fmt.Println("Database connection failed")
	}

	router := chi.NewRouter()

	router.Get("/", controllers.SaveArticles)
	router.Get("/articles", controllers.GetArticles)

	http.ListenAndServe(":8080", router)
}
