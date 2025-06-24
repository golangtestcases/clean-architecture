package main

import (
	"fmt"
	"os"

	"github.com/golangtestcases/clean-architecture/internal/app"
)

func main() {
	fmt.Println("app starting")

	app, err := app.NewApp(os.Getenv("CONFIG_ENV_VAR"))
	if err != nil {
		panic(err)
	}

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
}
