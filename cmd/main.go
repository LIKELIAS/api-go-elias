package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/elias/api-go-elias/internal/application"
	"github.com/elias/api-go-elias/internal/infrastructure/db"
	"github.com/elias/api-go-elias/internal/infrastructure/http/handler"
	"github.com/elias/api-go-elias/internal/infrastructure/http/middleware"
	"github.com/elias/api-go-elias/internal/infrastructure/repository"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: no se encontró archivo .env, usando variables del sistema")
	}

	// --- Infraestructura: conexión a base de datos ---
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Error al conectar con PostgreSQL: %v", err)
	}
	defer database.Close()
	log.Println("Conexión a PostgreSQL establecida correctamente")

	// --- Inyección de dependencias (Arquitectura Hexagonal) ---
	// 1. Repositorio (adaptador de salida)
	userRepo := repository.NewPostgresUserRepository(database)

	// 2. Servicio de aplicación (lógica de negocio)
	userSvc := application.NewUserService(userRepo)
	estudianteRepo := repository.NewPostgresEstudianteRepository(database)
	estudianteSvc := application.NewEstudianteService(estudianteRepo)
	estudianteHandler := handler.NewEstudianteHandler(estudianteSvc)

	// Archivo
	archivoRepo := repository.NewPostgresArchivoRepository(database)
	archivoSvc := application.NewArchivoService(archivoRepo)
	archivoHandler := handler.NewArchivoHandler(archivoSvc)

	// 3. Handler HTTP (adaptador de entrada)
	userHandler := handler.NewUserHandler(userSvc)

	// --- Router ---
	r := gin.Default()

	// Rutas públicas (sin autenticación)
	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// Rutas protegidas (requieren JWT)
	users := r.Group("/users")
	users.Use(middleware.JWTMiddleware())
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	estudiantes := r.Group("/estudiantes")
	estudiantes.Use(middleware.JWTMiddleware())
	{
		estudiantes.POST("", estudianteHandler.Crear)
		estudiantes.GET("", estudianteHandler.ObtenerTodos)
		estudiantes.GET("/buscar", estudianteHandler.BuscarPorFecha)
		estudiantes.GET("/:id", estudianteHandler.ObtenerPorID)
		estudiantes.PUT("/:id", estudianteHandler.Actualizar)
		estudiantes.DELETE("/:id", estudianteHandler.Eliminar)
	}

	archivos := r.Group("/archivos")
	archivos.Use(middleware.JWTMiddleware())
	{
		archivos.GET("", archivoHandler.Listar)
	}

	upload := r.Group("/upload")
	upload.Use(middleware.JWTMiddleware())
	{
		upload.POST("", archivoHandler.Subir)
	}

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
