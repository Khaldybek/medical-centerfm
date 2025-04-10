package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Get service URLs from environment variables
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	departmentServiceURL := getEnv("DEPARTMENT_SERVICE_URL", "http://localhost:8082")
	doctorServiceURL := getEnv("DOCTOR_SERVICE_URL", "http://localhost:8083")
	appointmentServiceURL := getEnv("APPOINTMENT_SERVICE_URL", "http://localhost:8084")

	// Auth routes
	router.POST("/api/v1/auth/register", func(c *gin.Context) {
		proxyRequest(c, authServiceURL+"/register")
	})

	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		proxyRequest(c, authServiceURL+"/login")
	})

	// Department routes
	router.POST("/api/v1/departments", func(c *gin.Context) {
		proxyRequest(c, departmentServiceURL+"/departments")
	})

	router.GET("/api/v1/departments", func(c *gin.Context) {
		proxyRequest(c, departmentServiceURL+"/departments")
	})

	router.GET("/api/v1/departments/:id", func(c *gin.Context) {
		proxyRequest(c, departmentServiceURL+"/departments/"+c.Param("id"))
	})

	router.PUT("/api/v1/departments/:id", func(c *gin.Context) {
		proxyRequest(c, departmentServiceURL+"/departments/"+c.Param("id"))
	})

	// Doctor routes
	router.POST("/api/v1/doctors", func(c *gin.Context) {
		proxyRequest(c, doctorServiceURL+"/doctors")
	})

	router.GET("/api/v1/doctors", func(c *gin.Context) {
		proxyRequest(c, doctorServiceURL+"/doctors")
	})

	router.GET("/api/v1/doctors/:id", func(c *gin.Context) {
		proxyRequest(c, doctorServiceURL+"/doctors/"+c.Param("id"))
	})

	router.PUT("/api/v1/doctors/:id", func(c *gin.Context) {
		proxyRequest(c, doctorServiceURL+"/doctors/"+c.Param("id"))
	})

	router.POST("/api/v1/doctors/:id/availability", func(c *gin.Context) {
		proxyRequest(c, doctorServiceURL+"/doctors/"+c.Param("id")+"/availability")
	})

	router.GET("/api/v1/doctors/:id/availability", func(c *gin.Context) {
		proxyRequest(c, doctorServiceURL+"/doctors/"+c.Param("id")+"/availability")
	})

	// Appointment routes
	router.POST("/api/v1/appointments", func(c *gin.Context) {
		proxyRequest(c, appointmentServiceURL+"/appointments")
	})

	router.GET("/api/v1/appointments", func(c *gin.Context) {
		proxyRequest(c, appointmentServiceURL+"/appointments")
	})

	router.GET("/api/v1/appointments/:id", func(c *gin.Context) {
		proxyRequest(c, appointmentServiceURL+"/appointments/"+c.Param("id"))
	})

	router.GET("/api/v1/appointments/doctor/:doctor_id", func(c *gin.Context) {
		proxyRequest(c, appointmentServiceURL+"/appointments/doctor/"+c.Param("doctor_id"))
	})

	router.GET("/api/v1/appointments/department/:department_id", func(c *gin.Context) {
		proxyRequest(c, appointmentServiceURL+"/appointments/department/"+c.Param("department_id"))
	})

	router.PUT("/api/v1/appointments/:id", func(c *gin.Context) {
		proxyRequest(c, appointmentServiceURL+"/appointments/"+c.Param("id"))
	})

	router.Run(":8080")
}

func proxyRequest(c *gin.Context, targetURL string) {
	// Create a new request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status code and body
	c.Status(resp.StatusCode)
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, resp.Body)
		return err == nil
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
