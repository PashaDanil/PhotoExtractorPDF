package main

import "imgpdf/internal/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	a.Run()
}
