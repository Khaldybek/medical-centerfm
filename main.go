package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	impl "medical-center/internal/gorm"
	"medical-center/internal/handler"
	"medical-center/internal/middleware"
	"medical-center/internal/migrations"
	"medical-center/internal/models/user"
	"medical-center/internal/service"
)

func main() {
	db, err := gorm.Open(postgres.Open("postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run migrations
	migrator := migrations.NewMigrator(db)
	migrator.AddMigration(&migrations.CreateDepartmentsTable{})
	migrator.AddMigration(&migrations.CreateDoctorsTable{})
	migrator.AddMigration(&migrations.CreateSchedulesTable{})
	migrator.AddMigration(&migrations.CreateAppointmentsTable{})
	migrator.AddMigration(&migrations.CreateUsersTable{})

	log.Println("Running database migrations...")
	if err := migrator.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Create default admin user if it doesn't exist
	createDefaultAdmin(db)

	deptRepo := impl.NewDepartmentRepository(db)
	doctorRepo := impl.NewDoctorRepository(db)
	scheduleRepo := impl.NewScheduleRepository(db)
	appointmentRepo := impl.NewAppoinmentRepository(db)
	userRepo := impl.NewUserRepository(db)

	deptService := service.NewDepartmentService(deptRepo)
	doctorService := service.NewDoctorService(doctorRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	authService := service.NewAuthService(userRepo)

	deptHandler := handler.NewDepartmentHandler(deptService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()

	// Public routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authService))
	{
		api.GET("/me", authHandler.Me)

		// Department routes - Admin only
		departments := api.Group("/departments")
		adminOnly := departments.Group("")
		adminOnly.Use(middleware.RoleMiddleware(user.RoleAdmin))
		{
			adminOnly.POST("", deptHandler.CreateDepartment)
			adminOnly.PUT("/:id", deptHandler.UpdateDepartment)
			//adminOnly.DELETE("/:id", deptHandler.DeleteDepartment)
		}
		// Public department routes (still require authentication)
		departments.GET("", deptHandler.GetAllDepartments)
		departments.GET("/:id", deptHandler.GetDepartment)
		departments.GET("/:id/slots", deptHandler.GetDepartmentSlots)

		// Doctor routes
		doctors := api.Group("/doctors")
		doctorAdmin := doctors.Group("")
		doctorAdmin.Use(middleware.RoleMiddleware(user.RoleAdmin, user.RoleDoctor))
		{
			doctorAdmin.PATCH("/:id/availability", doctorHandler.SetAvailability)
		}
		adminOnly = doctors.Group("")
		adminOnly.Use(middleware.RoleMiddleware(user.RoleAdmin))
		{
			adminOnly.POST("", doctorHandler.CreateDoctor)
			adminOnly.PUT("/:id", doctorHandler.UpdateDoctor)
			//adminOnly.DELETE("/:id", doctorHandler.DeleteDoctor)
		}
		doctors.GET("", doctorHandler.GetAllDoctors)
		doctors.GET("/:id", doctorHandler.GetDoctor)

		// Schedule routes
		schedules := api.Group("/schedules")
		scheduleAdmin := schedules.Group("")
		scheduleAdmin.Use(middleware.RoleMiddleware(user.RoleAdmin, user.RoleDoctor))
		{
			scheduleAdmin.POST("", scheduleHandler.CreateSlot)
		}
		schedules.GET("/:id", scheduleHandler.GetSlot)
		schedules.GET("/doctor/:doctor_id", scheduleHandler.GetDoctorSlots)
		schedules.POST("/:id/book", scheduleHandler.BookSlot)
		schedules.GET("/available", scheduleHandler.GetAvailableSlots)

		// Appointment routes
		appointments := api.Group("/appointments")
		appointmentAdmin := appointments.Group("")
		appointmentAdmin.Use(middleware.RoleMiddleware(user.RoleAdmin, user.RoleDoctor))
		{
			appointmentAdmin.PUT("/:id", appointmentHandler.UpdateAppointment)
			//appointmentAdmin.DELETE("/:id", appointmentHandler.DeleteAppointment)
		}
		appointments.POST("", appointmentHandler.CreateAppointment)
		appointments.GET("", appointmentHandler.GetAllAppointments)
		appointments.GET("/:id", appointmentHandler.GetAppointment)
		appointments.GET("/department/:department_id", appointmentHandler.GetAppointmentsByDepartment)
	}

	router.Run(":8080")
}

func createDefaultAdmin(db *gorm.DB) {
	var adminUser user.User
	result := db.Where("email = ?", "admin@example.com").First(&adminUser)
	if result.RowsAffected == 0 {
		// Create admin user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser = user.User{
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			Name:     "Admin User",
			Role:     user.RoleAdmin,
		}
		db.Create(&adminUser)
	}
}
