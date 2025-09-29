package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"user_service/models"
	"user_service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c echo.Context) error {
	var request models.RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to hash password"})
	}
	request.Password = string(hashedPassword)

	err = repositories.CreateUser(&models.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User registered successfully"})
}

func LoginUser(c echo.Context) error {
	var request models.LoginRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	user, err := repositories.GetUserByEmail(request.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid password"})
	}

	if user.Role != "user" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid role"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("USER_JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to sign token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Login successful", "token": tokenString})
}

func RegisterAdmin(c echo.Context) error {
	var request models.RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to hash password"})
	}
	request.Password = string(hashedPassword)

	err = repositories.CreateUser(&models.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User registered successfully"})
}

func LoginAdmin(c echo.Context) error {
	var request models.LoginRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	user, err := repositories.GetUserByEmail(request.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid password"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid role"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("ADMIN_JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to sign token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Login successful", "token": tokenString})
}

func GetAllUsers(c echo.Context) error {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get users"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Users fetched successfully", "users": users})
}

func GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid ID"})
	}

	user, err := repositories.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User fetched successfully", "user": user})
}

func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid ID"})
	}

	var request models.UpdateUserRequest
	err = c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	err = repositories.UpdateUser(id, &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User updated successfully"})
}

func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid ID"})
	}

	err = repositories.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}
