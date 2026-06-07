package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/elias/api-go-elias/internal/application"
	"github.com/elias/api-go-elias/internal/infrastructure/db"
	"github.com/elias/api-go-elias/internal/infrastructure/http/handler"
	"github.com/elias/api-go-elias/internal/infrastructure/http/middleware"
	"github.com/elias/api-go-elias/internal/infrastructure/repository"
)

var ginLambda *ginadapter.GinLambdaV2

func setupRouter() *gin.Engine {
	// --- Infraestructura: conexión a base de datos ---
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Error al conectar con PostgreSQL: %v", err)
	}
	log.Println("Conexión a PostgreSQL establecida correctamente")

	// --- Inyección de dependencias (Arquitectura Hexagonal) ---
	userRepo := repository.NewPostgresUserRepository(database)
	userSvc := application.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	estudianteRepo := repository.NewPostgresEstudianteRepository(database)
	estudianteSvc := application.NewEstudianteService(estudianteRepo)
	estudianteHandler := handler.NewEstudianteHandler(estudianteSvc)

	archivoRepo := repository.NewPostgresArchivoRepository(database)
	archivoSvc := application.NewArchivoService(archivoRepo)
	archivoHandler := handler.NewArchivoHandler(archivoSvc)

	// --- Router ---
	r := gin.New()
	r.Use(gin.Recovery())

	// Health check (útil para verificar el deploy en AWS)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-go-elias",
		})
	})

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

	return r
}

func handler_lambda(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Cargar variables de entorno (solo aplica en modo local)
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: no se encontró archivo .env, usando variables del sistema")
	}

	r := setupRouter()

	// Si corre en Lambda, LAMBDA_TASK_ROOT viene definido por AWS
	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		// Modo Lambda
		log.Println("Iniciando en modo AWS Lambda")
		gin.SetMode(gin.ReleaseMode)
		ginLambda = ginadapter.NewV2(r)
		lambda.Start(handler_lambda)
	} else {
		// Modo local (docker, go run, etc.)
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Servidor iniciado en http://localhost:%s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}
}
