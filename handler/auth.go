package handler

import (
	"auth_project/cmd/jwt_auth/ent"
	"auth_project/cmd/jwt_auth/ent/user"
	"context"
	"fmt"
	"log"
	"time"

	"auth_project/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RegisterUser struct {
			Name     string `json:"name"`
			LastName string `json:"lastname"`
			UserName string `json:"username"`
			Email    string `json:"email"`
			Pass     string `json:"pass"`
			Role     string `json:"role"`
		}

		user := new(RegisterUser)
		if err := c.BodyParser(user); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "errors": err.Error()})
		}

		hash, err := hashPassword(user.Pass)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to hash password", "errors": err.Error()})
		}

		if user.Role != "user" && user.Role != "admin" {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid role"})
		}

		user.Pass = hash
		u, err := client.User.
			Create().
			SetName(user.Name).
			SetLastname(user.LastName).
			SetUsername(user.UserName).
			SetEmail(user.Email).
			SetPassword(user.Pass).
			SetRole(user.Role).
			Save(context.Background())

		if err != nil {
			return fmt.Errorf("failed creating user: %w", err)
		}

		log.Println("user was created: ", u)
		return nil
	}
}

func Login(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginUser struct {
			UserName string `json:"username"`
			Pass     string `json:"password"`
		}

		loginUser := new(LoginUser)

		if err := c.BodyParser(loginUser); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"errors":  err.Error(),
			})
		}

		u, err := client.User.
			Query().
			Where(user.Username(loginUser.UserName)).
			Only(context.Background())

		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
				"errors":  err.Error()})
		}

		if !CheckPasswordHash(loginUser.Pass, u.Password) {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid credentials",
			})
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = u.Username
		claims["user_id"] = u.ID
		claims["role"] = u.Role
		claims["exp"] = time.Now().Add(time.Minute * 20).Unix()

		t, err := token.SignedString([]byte(config.LoadEnvVariable("SECRET")))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		log.Println("User logged in successfully: ", u)
		return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
	}
}

func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "You are authorized"})
}
