package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/elias/api-go-elias/internal/domain"
)

// UserHandler maneja las peticiones HTTP para el recurso usuario
type UserHandler struct {
	service domain.UserService
}

// NewUserHandler crea una nueva instancia del handler
func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// --- Request DTOs ---

type registerRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type updateRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// --- Handlers ---

// Register godoc
// @Summary  Registrar un nuevo usuario
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body registerRequest true "Datos del usuario"
// @Success  201 {object} domain.User
// @Router   /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary  Iniciar sesión y obtener JWT
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body loginRequest true "Credenciales"
// @Success  200 {object} map[string]string
// @Router   /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetAll godoc
// @Summary  Obtener todos los usuarios
// @Tags     users
// @Security BearerAuth
// @Produce  json
// @Success  200 {array} domain.User
// @Router   /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetByID godoc
// @Summary  Obtener un usuario por ID
// @Tags     users
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "ID del usuario"
// @Success  200 {object} domain.User
// @Router   /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary  Actualizar un usuario
// @Tags     users
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int           true "ID del usuario"
// @Param    body body updateRequest true "Nuevos datos"
// @Success  200 {object} domain.User
// @Router   /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary  Eliminar un usuario
// @Tags     users
// @Security BearerAuth
// @Param    id path int true "ID del usuario"
// @Success  204
// @Router   /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// parseID extrae y valida el parámetro :id de la URL
func parseID(c *gin.Context) (uint, error) {
	raw, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return 0, err
	}
	return uint(raw), nil
}

func (h *EstudianteHandler) BuscarPorFecha(c *gin.Context) {
	desde, err := time.Parse("2006-01-02", c.Query("desde"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'desde' inválido, use YYYY-MM-DD"})
		return
	}
	hasta, err := time.Parse("2006-01-02", c.Query("hasta"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'hasta' inválido, use YYYY-MM-DD"})
		return
	}
	lista, err := h.service.BuscarPorFechaIncripcion(desde, hasta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lista)
}
