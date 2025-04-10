package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"medical-center/internal/service"
)

type DoctorHandler struct {
	service *service.DoctorService
}

func NewDoctorHandler(s *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{service: s}
}

func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var request struct {
		Name         string `json:"name"`
		DepartmentID uint   `json:"department_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	doctor, err := h.service.CreateDoctor(request.Name, request.DepartmentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doctor)
}

func (h *DoctorHandler) GetDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doctor, err := h.service.GetDoctorByID(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var request struct {
		Name         string `json:"name"`
		DepartmentID uint   `json:"department_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	doctor, err := h.service.UpdateDoctor(uint(id), request.Name, request.DepartmentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) GetAllDoctors(c *gin.Context) {
	doctors, err := h.service.GetAllDoctors()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

func (h *DoctorHandler) SetAvailability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var request struct {
		Available bool `json:"available"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.SetAvailability(uint(id), request.Available); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
