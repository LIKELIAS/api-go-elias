package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/elias/api-go-elias/internal/domain"
)

type EstudianteHandler struct {
	service domain.EstudianteService
}

func NewEstudianteHandler(service domain.EstudianteService) *EstudianteHandler {
	return &EstudianteHandler{service: service}
}

type crearEstudianteRequest struct {
	Nombre          string `json:"nombre"           binding:"required"`
	ApellidoPat     string `json:"apellido_pat"     binding:"required"`
	ApellidoMat     string `json:"apellido_mat"     binding:"required"`
	FechaNacimiento string `json:"fecha_nacimiento" binding:"required"`
	Matricula       string `json:"matricula"        binding:"required"`
	Carrera         string `json:"carrera"          binding:"required"`
	FechaIncripcion string `json:"fecha_incripcion" binding:"required"`
}

type actualizarEstudianteRequest struct {
	Nombre      string `json:"nombre"       binding:"required"`
	ApellidoPat string `json:"apellido_pat" binding:"required"`
	ApellidoMat string `json:"apellido_mat" binding:"required"`
	Matricula   string `json:"matricula"    binding:"required"`
	Carrera     string `json:"carrera"      binding:"required"`
}

func (h *EstudianteHandler) Crear(c *gin.Context) {
	var req crearEstudianteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fechaNac, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha_nacimiento inválida, use formato YYYY-MM-DD"})
		return
	}

	fechaInc, err := time.Parse("2006-01-02", req.FechaIncripcion)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha_incripcion inválida, use formato YYYY-MM-DD"})
		return
	}

	e, err := h.service.CrearEstudiante(req.Nombre, req.ApellidoPat, req.ApellidoMat, req.Matricula, req.Carrera, fechaNac, fechaInc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *EstudianteHandler) ObtenerTodos(c *gin.Context) {
	lista, err := h.service.ObtenerTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lista)
}

func (h *EstudianteHandler) ObtenerPorID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	e, err := h.service.ObtenerEstudiante(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *EstudianteHandler) Actualizar(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	var req actualizarEstudianteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := h.service.ActualizarEstudiante(id, req.Nombre, req.ApellidoPat, req.ApellidoMat, req.Matricula, req.Carrera)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *EstudianteHandler) Eliminar(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	if err := h.service.EliminarEstudiante(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
