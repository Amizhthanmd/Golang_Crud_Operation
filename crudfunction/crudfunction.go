package crudfunction

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"Crud_operation_go/models"
)

var (
	db *gorm.DB
	updateEmailWG sync.WaitGroup
)

func SetupDB() {
	dsn := "host=localhost user=postgres password=Ami160320! dbname=Student port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&models.User{})
}

func CreateUser(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

func GetUser(c *gin.Context) {

	var user models.User

	userID := c.Param("id")

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetAllUsers(c *gin.Context) {

	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func UpdateUser(c *gin.Context) {

	var user models.User

	userID := c.Param("id")

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

//-------

func UpdateEmails() {
	for {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			log.Println("Failed to fetch users:", err)
		}

		for _, user := range users {
			updateEmailWG.Add(1)

			go func(u models.User) {
				defer updateEmailWG.Done()

				email := GenerateRandomEmail()

				u.Email = email
				if err := db.Save(&u).Error; err != nil {
					log.Printf("Failed to update email for user %d: %v", u.ID, err)
				}

			}(user)
		}

		updateEmailWG.Wait()

		time.Sleep(10 * time.Second)
	}
}

func GenerateRandomEmail() string {
	
	const randomlet = "abcdefghijklmnopqrstuvwxyz0123456789"

	result := make([]byte, 8)
	for i := range result {
		result[i] = randomlet[rand.Intn(len(randomlet))]
	}
	var domain = "gmail.com"

	email := string(result) + "@" + domain

	return email
}

func DeleteUser(c *gin.Context) {

	var user models.User
	userID := c.Param("id")

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
