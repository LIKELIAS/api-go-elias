package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
		log.Printf("Error al obtener archivo del form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se proporcionó archivo"})
		return
	}

	log.Printf("Archivo recibido: %s", file.Filename)

	// Abrir el archivo en memoria
	src, err := file.Open()
	if err != nil {
		log.Printf("Error al abrir archivo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al leer archivo"})
		return
	}
	defer src.Close()

	// Generar nombre único
	ext := filepath.Ext(file.Filename)
	nuevoNombre := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Obtener nombre del bucket S3 desde variable de entorno
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("AWS_REGION")
	log.Printf("S3_BUCKET: '%s', AWS_REGION: '%s'", bucket, region)

	if region == "" {
		region = "us-east-2"
	}

	var ruta string

	if bucket != "" {
		// Modo producción: subir a S3
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			log.Printf("Error al configurar AWS: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al configurar AWS"})
			return
		}

		client := s3.NewFromConfig(cfg)
		key := "uploads/" + nuevoNombre

		contentType := file.Header.Get("Content-Type")
		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      &bucket,
			Key:         &key,
			Body:        src,
			ContentType: &contentType,
		})
		if err != nil {
			log.Printf("Error al subir a S3: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al subir archivo a S3"})
			return
		}

		ruta = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
		log.Printf("Archivo subido a S3: %s", ruta)
	} else {
		// Modo local: guardar en disco
		ruta = "uploads/" + nuevoNombre
		if err := c.SaveUploadedFile(file, ruta); err != nil {
			log.Printf("Error al guardar en disco: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar archivo"})
			return
		}
	}

	userID, _ := c.Get("user_id")
	uid := uint(userID.(float64))

	archivo, err := h.service.SubirArchivo(file.Filename, ext, ruta, uid)
	if err != nil {
		log.Printf("Error al guardar en BD: %v", err)
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
