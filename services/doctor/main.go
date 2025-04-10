package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Name           string         `json:"name" binding:"required"`
	Email          string         `json:"email" binding:"required,email" gorm:"unique"`
	DepartmentID   uint           `json:"department_id" binding:"required"`
	Specialization string         `json:"specialization"`
	Availability   []Availability `json:"availability"`
}

type Availability struct {
	gorm.Model
	DoctorID  uint   `json:"doctor_id"`
	Day       string `json:"day" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
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
	db.AutoMigrate(&Doctor{}, &Availability{})

	router := gin.Default()

	// Doctor routes
	router.POST("/doctors", func(c *gin.Context) {
		var doctor Doctor
		if err := c.ShouldBindJSON(&doctor); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&doctor).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create doctor"})
			return
		}

		c.JSON(201, doctor)
	})

	router.GET("/doctors", func(c *gin.Context) {
		var doctors []Doctor
		if err := db.Preload("Availability").Find(&doctors).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch doctors"})
			return
		}

		c.JSON(200, doctors)
	})

	router.GET("/doctors/:id", func(c *gin.Context) {
		var doctor Doctor
		if err := db.Preload("Availability").First(&doctor, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Doctor not found"})
			return
		}

		c.JSON(200, doctor)
	})

	router.PUT("/doctors/:id", func(c *gin.Context) {
		var doctor Doctor
		if err := db.First(&doctor, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Doctor not found"})
			return
		}

		if err := c.ShouldBindJSON(&doctor); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&doctor).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update doctor"})
			return
		}

		c.JSON(200, doctor)
	})

	// Availability routes
	router.POST("/doctors/:id/availability", func(c *gin.Context) {
		var availability Availability
		if err := c.ShouldBindJSON(&availability); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		availability.DoctorID = uint(c.GetUint("id"))
		if err := db.Create(&availability).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create availability"})
			return
		}

		c.JSON(201, availability)
	})

	router.GET("/doctors/:id/availability", func(c *gin.Context) {
		var availabilities []Availability
		if err := db.Where("doctor_id = ?", c.Param("id")).Find(&availabilities).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch availabilities"})
			return
		}

		c.JSON(200, availabilities)
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
