package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/elias/api-go-elias/internal/domain"
	"github.com/gin-gonic/gin"
)

type ArchivoHandler struct {
	service domain.ArchivoService
}

func NewArchivoHandler(service domain.ArchivoService) *ArchivoHandler {
	return &ArchivoHandler{service: service}
}

func (h *ArchivoHandler) Subir(c *gin.Context) {
	file, err := c.FormFile("archivo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se proporcionó archivo"})
		return
	}

	ext := filepath.Ext(file.Filename)
	nuevoNombre := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	ruta := "uploads/" + nuevoNombre

	if err := c.SaveUploadedFile(file, ruta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar archivo"})
		return
	}

	userID, _ := c.Get("user_id")
	uid := uint(userID.(float64))

	archivo, err := h.service.SubirArchivo(file.Filename, ext, ruta, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, archivo)
}

func (h *ArchivoHandler) Listar(c *gin.Context) {
	archivos, err := h.service.GetArchivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, archivos)
}
