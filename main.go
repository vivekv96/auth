package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/vivekv96/auth/database"
	"github.com/vivekv96/auth/routes"
)

func main() {
	if err := database.ConnectToMySQL(&database.MySQLConfig{
		Host:     "mysql-server",
		Username: "root",
		Password: "root123",
		Port:     3306,
		DBName:   "mysql",
	}); err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if err := app.Listen(":8000"); err != nil {
		log.Fatalln(err)
	}
}
