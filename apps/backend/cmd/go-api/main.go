package main

import (
	"context"
	_ "go-api/docs"
	"go-api/internal/app"
	"log"
)

// @title PDF to Images API
// @version 1.0
// @description API for uploading PDF files and extracting images from them
// @description
// @description This API provides endpoints for:
// @description - Uploading PDF files directly or via presigned URLs
// @description - Processing PDFs to extract images
// @description - Downloading extracted images as ZIP archives
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	ctx := context.Background()
	a, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
