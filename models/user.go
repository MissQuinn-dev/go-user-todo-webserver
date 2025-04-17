package models

import (
	"app/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string  `json:"name"`
	TotalPoints int     `json:"totalPoints" gorm:"default:0"`
	Tasks       *[]Task `json:"tasks"`
}

func (u *User) AddPoints(tx *gorm.DB, points int) error {
	return tx.Model(u).UpdateColumn("total_points", gorm.Expr("total_points + ?", points)).Error
}

func NewUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(201).SendString(err.Error())
	}

	db.Create(&user)
	return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var user User
	db.Find(&user, id)
	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var users []User
	db.Preload("Tasks").Find(&users)
	return c.JSON(users)
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&user)
	return c.Status(200).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var user User
	db.First(&user, id)
	if user.Name == "" {
		return c.Status(500).SendString("No user was found with this ID")
	}

	db.Delete(&user)
	return c.SendString("User was Removed")
}
