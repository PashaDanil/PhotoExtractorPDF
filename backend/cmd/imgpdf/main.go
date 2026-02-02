package main

import (
	"context"
	_ "imgpdf/docs"
	"imgpdf/internal/app"
	"log"
)

// @title PDF Image Extractor API
// @version 1.0
// @description API для загрузки PDF и получения ZIP с изображениями
// @host localhost:8080
// @BasePath /
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
