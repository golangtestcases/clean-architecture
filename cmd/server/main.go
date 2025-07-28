// @title Subscription Service API
// @version 1.0
// @description REST-сервис для агрегации данных об онлайн-подписках пользователей
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"os"

	"github.com/golangtestcases/subscribe-service/internal/app"
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
