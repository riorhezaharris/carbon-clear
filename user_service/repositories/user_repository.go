package repositories

import (
	"time"
	"user_service/configs"
	"user_service/models"
)

func CreateUser(payload *models.User) error {
	err := configs.DB.Create(payload).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := configs.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := configs.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := configs.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id int, payload *models.User) error {
	var user models.User
	err := configs.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = payload.Password
	user.UpdatedAt = time.Now()

	err = configs.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id int) error {
	var user models.User
	err := configs.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	err = configs.DB.Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
