variable "aws_region" {
  description = "Región de AWS donde se desplegará la infraestructura"
  type        = string
  default     = "us-east-2"
}

variable "project_name" {
  description = "Nombre del proyecto (usado para nombrar recursos)"
  type        = string
  default     = "api-go-elias"
}

variable "environment" {
  description = "Ambiente de despliegue"
  type        = string
  default     = "production"
}

variable "lambda_memory_mb" {
  description = "Memoria asignada a la función Lambda en MB"
  type        = number
  default     = 256
}

variable "lambda_timeout_seconds" {
  description = "Tiempo máximo de ejecución de Lambda en segundos"
  type        = number
  default     = 30
}

variable "database_url" {
  description = "URL de conexión a la base de datos (Neon PostgreSQL)"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "Clave secreta para firmar tokens JWT"
  type        = string
  sensitive   = true
}

variable "lambda_zip_path" {
  description = "Ruta local al archivo ZIP del Lambda"
  type        = string
  default     = "../lambda.zip"
}
