package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "myuser")
	dbPassword := getEnv("DB_PASSWORD", "mypassword")
	dbName := getEnv("DB_NAME", "mydatabase")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Department{})

	router := gin.Default()

	// Department routes
	router.POST("/departments", func(c *gin.Context) {
		var dept Department
		if err := c.ShouldBindJSON(&dept); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&dept).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create department"})
			return
		}

		c.JSON(201, dept)
	})

	router.GET("/departments", func(c *gin.Context) {
		var departments []Department
		if err := db.Find(&departments).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch departments"})
			return
		}

		c.JSON(200, departments)
	})

	router.GET("/departments/:id", func(c *gin.Context) {
		var dept Department
		if err := db.First(&dept, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Department not found"})
			return
		}

		c.JSON(200, dept)
	})

	router.PUT("/departments/:id", func(c *gin.Context) {
		var dept Department
		if err := db.First(&dept, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Department not found"})
			return
		}

		if err := c.ShouldBindJSON(&dept); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&dept).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update department"})
			return
		}

		c.JSON(200, dept)
	})

	router.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
