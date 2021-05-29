package handlers

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vivekv96/auth/database"
	"github.com/vivekv96/auth/models"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var urlFormat = "Click <a href=\"%s\">here</a> to reset your password."

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	token := generateRandString(12)

	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	database.Gorm.Create(&passwordReset)

	from := "admin@auth-server.com"
	to := []string{data["email"]}
	url := "http://localhost:3000/reset/" + token
	message := []byte(fmt.Sprintf(urlFormat, url))

	// MailHog SMTP bind address
	err := smtp.SendMail("0.0.0.0:1025", nil, from, to, message)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func generateRandString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

func Reset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// return error message if `password` & `confirmPassword` fields are not identical
	if data["password"] != data["confirmPassword"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "passwords do not match!",
		})
	}

	var passwordReset models.PasswordReset

	if err := database.Gorm.Where("token = ?", data["token"]).Last(&passwordReset).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

	database.Gorm.Model(&models.User{}).Where("email = ?", passwordReset.Email).
		Update("password", password)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
