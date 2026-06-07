# ─────────────────────────────────────────────────────────────────────────────
# terraform.tfvars  —  NO commitear este archivo con valores reales.
# Las variables sensibles (database_url, jwt_secret) se inyectan desde
# GitHub Actions Secrets como variables de entorno TF_VAR_*.
# ─────────────────────────────────────────────────────────────────────────────

aws_region   = "us-east-2"
project_name = "api-go-elias"
environment  = "production"

lambda_memory_mb       = 256
lambda_timeout_seconds = 30