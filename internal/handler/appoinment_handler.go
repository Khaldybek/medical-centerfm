package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"medical-center/internal/service"
)

type AppointmentHandler struct {
	service *service.AppointmentService
}

func NewAppointmentHandler(s *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: s}
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var request struct {
		PatientName     string    `json:"patient_name"`
		Email           string    `json:"email"`
		Phone           string    `json:"phone"`
		DepartmentID    uint      `json:"department_id"`
		DoctorID        uint      `json:"doctor_id"`
		AppointmentTime time.Time `json:"appointment_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	appt, err := h.service.CreateAppointment(
		request.PatientName,
		request.Email,
		request.Phone,
		request.DepartmentID,
		request.DoctorID,
		request.AppointmentTime,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appt)
}

func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appt, err := h.service.GetAppointmentByID(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var request struct {
		PatientName     string    `json:"patient_name"`
		Email           string    `json:"email"`
		Phone           string    `json:"phone"`
		DepartmentID    uint      `json:"department_id"`
		DoctorID        uint      `json:"doctor_id"`
		AppointmentTime time.Time `json:"appointment_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	appt, err := h.service.UpdateAppointment(
		uint(id),
		request.PatientName,
		request.Email,
		request.Phone,
		request.DepartmentID,
		request.DoctorID,
		request.AppointmentTime,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) GetAllAppointments(c *gin.Context) {
	appointments, err := h.service.GetAllAppointments()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentHandler) GetAppointmentsByDepartment(c *gin.Context) {
	deptIDStr := c.Param("department_id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	appointments, err := h.service.GetByDepartment(uint(deptID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}
