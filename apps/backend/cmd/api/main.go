package main

import (
	_ "api/docs"
	"api/internal/app"
	"log"
)

// @title PDF to Images API
// @version 1.0
// @description API for converting PDF documents to images
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}
