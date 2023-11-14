package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"type:varchar(100);not null"`
	Email  string `gorm:"type:varchar(100);unique;not null"`
	Passwd string `gorm:"not null"`
	Active bool   `gorm:"type:boolean;not null;default:true"`
}

func (User) TableName() string {
	return "users"
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Passwd), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Replace the plain text password with the hashed password
	newUser.Passwd = string(hashedPassword)

	// Create a new user record in the database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user created", "user": newUser})
}

func getUsers(c *gin.Context) {
	var users []User

	// Retrieve all users from the database, excluding the "Passwd" field
	if err := db.Select("id, name, email, active").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func deleteUser(c *gin.Context) {
	// Get the user's ID from the request parameters
	userID := c.Param("id")

	// Check if the user with the given ID exists
	var user User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user from the database
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
