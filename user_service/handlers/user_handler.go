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

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User registration details"
// @Success 200 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /api/users/register [post]
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

// LoginUser godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with token"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid credentials or role"
// @Failure 500 {object} map[string]string "Failed to get user or sign token"
// @Router /api/users/login [post]
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

// RegisterAdmin godoc
// @Summary Register a new admin
// @Description Register a new admin account
// @Tags admin
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Admin registration details"
// @Success 200 {object} map[string]string "Admin registered successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to create admin"
// @Router /admin/users/register [post]
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

// LoginAdmin godoc
// @Summary Login admin
// @Description Authenticate admin and return JWT token
// @Tags admin
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with token"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid credentials or role"
// @Failure 500 {object} map[string]string "Failed to get user or sign token"
// @Router /admin/users/login [post]
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Success 200 {object} map[string]interface{} "Users fetched successfully"
// @Failure 500 {object} map[string]string "Failed to get users"
// @Router /admin/users [get]
func GetAllUsers(c echo.Context) error {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get users"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Users fetched successfully", "users": users})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User fetched successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Failed to get user"
// @Router /admin/users/{id} [get]
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

// UpdateUser godoc
// @Summary Update user
// @Description Update a user's information (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "User update details"
// @Success 200 {object} map[string]string "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid ID or request body"
// @Failure 500 {object} map[string]string "Failed to update user"
// @Router /admin/users/{id} [put]
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

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Failed to delete user"
// @Router /admin/users/{id} [delete]
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

// GetProfile godoc
// @Summary Get user profile
// @Description Get the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security UserAuth
// @Success 200 {object} map[string]interface{} "Profile fetched successfully"
// @Failure 500 {object} map[string]string "Failed to get user profile"
// @Router /api/users/profile [get]
func GetProfile(c echo.Context) error {
	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// Get user from database
	userData, err := repositories.GetUserByID(int(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user profile"})
	}

	// Remove password from response
	userData.Password = ""

	return c.JSON(http.StatusOK, echo.Map{"message": "Profile fetched successfully", "user": userData})
}
