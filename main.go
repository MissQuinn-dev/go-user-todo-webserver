package main

import (
	"app/database"
	"app/models"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	// Task
	app.Post("/api/v1/task", models.NewTask)
	app.Get("/api/v1/task/:id", models.GetTask)
	app.Get("/api/v1/task", models.GetTasks)
	app.Put("/api/v1/task/:id", models.UpdateTask)
	app.Delete("/api/v1/task/:id", models.DeleteTask)

	//User
	app.Post("/api/v1/user", models.NewUser)
	app.Get("/api/v1/user/:id", models.GetUser)
	app.Get("/api/v1/user", models.GetUsers)
	app.Put("/api/v1/user/:id", models.UpdateUser)
	app.Delete("/api/v1/user/:id", models.DeleteUser)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("models.db"))
	if err != nil {
		panic("sowwie >.< the database connection failed")
	}

	fmt.Println("Conection Opened to Database")
	database.DBConn.AutoMigrate(&models.User{}, &models.Task{})
	fmt.Println("Database Migrated")
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	initDatabase()

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
