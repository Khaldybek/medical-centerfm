package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "myuser")
	dbPassword := getEnv("DB_PASSWORD", "mypassword")
	dbName := getEnv("DB_NAME", "mydatabase")
	// JWT Secret available for future JWT token implementation
	_ = getEnv("JWT_SECRET", "your-secret-key")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&User{})

	// Create default admin if not exists
	createDefaultAdmin(db)

	router := gin.Default()

	// Auth routes
	router.POST("/register", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
			Name     string `json:"name" binding:"required"`
			Role     string `json:"role" binding:"required,oneof=admin doctor"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Check if user already exists
		var existingUser User
		if result := db.Where("email = ?", input.Email).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(400, gin.H{"error": "User already exists"})
			return
		}

		// Create new user
		user := User{
			Email:    input.Email,
			Password: input.Password, // In production, hash the password
			Name:     input.Name,
			Role:     input.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(201, gin.H{"message": "User created successfully"})
	})

	router.POST("/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var user User
		if result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&user); result.Error != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// In production, generate JWT token here
		c.JSON(200, gin.H{
			"message": "Login successful",
			"user": gin.H{
				"id":    user.Model.ID,
				"email": user.Email,
				"name":  user.Name,
				"role":  user.Role,
			},
		})
	})

	router.Run(":8080")
}

func createDefaultAdmin(db *gorm.DB) {
	var admin User
	if result := db.Where("email = ?", "admin@example.com").First(&admin); result.RowsAffected == 0 {
		admin = User{
			Email:    "admin@example.com",
			Password: "admin123", // In production, use hashed password
			Name:     "Admin User",
			Role:     "admin",
		}
		db.Create(&admin)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
