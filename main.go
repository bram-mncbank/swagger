package main

import (
	"swaggo/routers"
)

// @title Gin Swagger Simple API
// @version 1.0
// @description This is a simple API server for training.
// @termsOfService http://swagger.io/terms/

// @contact.name Brahmantyo
// @contact.url http://www.mncbank.co.id/
// @contact.email brahmantyo.adi@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8181
// @BasePath /
// @schemes http
func main() {
	PORT := ":8181"
	server := routers.StartServer()

	server.Run(PORT)
}
