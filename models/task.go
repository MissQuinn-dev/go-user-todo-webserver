package models

import (
	"app/database"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name   string `json:"name"`
	Points int    `json:"points"`
	UserID *uint  `json:"user_id"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func NewTask(c *fiber.Ctx) error {
	db := database.DBConn
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&task)
	return c.JSON(task)
}

func GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var task Task
	db.Find(&task, id)
	return c.JSON(task)
}

func GetTasks(c *fiber.Ctx) error {
	db := database.DBConn
	var tasks []Task
	db.Find(&tasks)
	return c.JSON(tasks)
}

func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	task := new(Task)

	if err := c.BodyParser(task); err != nil {
		return c.Status(201).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&task)
	return c.Status(200).JSON(task)
}

func (t *Task) BeforeDelete(tx *gorm.DB) (err error) {
	if t.UserID == nil {
		return nil
	}
	var user User
	err = tx.First(&user, *t.UserID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return user.AddPoints(tx, t.Points)
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var task Task
	db.First(&task, id)

	if task.Name == "" {
		return c.Status(500).SendString("No task fround with ID")
	}

	db.Delete(&task)
	return c.SendString("Task Completed and Removed from Database")
}
