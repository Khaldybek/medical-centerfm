package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	PatientName     string    `json:"patient_name" binding:"required"`
	PatientPhone    string    `json:"patient_phone" binding:"required"`
	DoctorID        uint      `json:"doctor_id" binding:"required"`
	DepartmentID    uint      `json:"department_id" binding:"required"`
	AppointmentTime time.Time `json:"appointment_time" binding:"required"`
	Status          string    `json:"status" gorm:"default:'scheduled'"`
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
	db.AutoMigrate(&Appointment{})

	router := gin.Default()

	// Appointment routes
	router.POST("/appointments", func(c *gin.Context) {
		var appointment Appointment
		if err := c.ShouldBindJSON(&appointment); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Check if the time slot is available
		var existingAppointment Appointment
		if err := db.Where("doctor_id = ? AND appointment_time = ?",
			appointment.DoctorID, appointment.AppointmentTime).First(&existingAppointment).Error; err == nil {
			c.JSON(400, gin.H{"error": "This time slot is already booked"})
			return
		}

		if err := db.Create(&appointment).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create appointment"})
			return
		}

		c.JSON(201, appointment)
	})

	router.GET("/appointments", func(c *gin.Context) {
		var appointments []Appointment
		if err := db.Find(&appointments).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch appointments"})
			return
		}

		c.JSON(200, appointments)
	})

	router.GET("/appointments/:id", func(c *gin.Context) {
		var appointment Appointment
		if err := db.First(&appointment, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Appointment not found"})
			return
		}

		c.JSON(200, appointment)
	})

	router.GET("/appointments/doctor/:doctor_id", func(c *gin.Context) {
		var appointments []Appointment
		if err := db.Where("doctor_id = ?", c.Param("doctor_id")).Find(&appointments).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch appointments"})
			return
		}

		c.JSON(200, appointments)
	})

	router.GET("/appointments/department/:department_id", func(c *gin.Context) {
		var appointments []Appointment
		if err := db.Where("department_id = ?", c.Param("department_id")).Find(&appointments).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch appointments"})
			return
		}

		c.JSON(200, appointments)
	})

	router.PUT("/appointments/:id", func(c *gin.Context) {
		var appointment Appointment
		if err := db.First(&appointment, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Appointment not found"})
			return
		}

		if err := c.ShouldBindJSON(&appointment); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&appointment).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update appointment"})
			return
		}

		c.JSON(200, appointment)
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
