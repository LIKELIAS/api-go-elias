output "api_gateway_url" {
  description = "URL base del API Gateway — úsala en la app Android"
  value       = aws_apigatewayv2_api.http_api.api_endpoint
}

output "lambda_function_name" {
  description = "Nombre de la función Lambda desplegada"
  value       = aws_lambda_function.api.function_name
}

output "lambda_arn" {
  description = "ARN de la función Lambda"
  value       = aws_lambda_function.api.arn
}

output "s3_bucket_name" {
  description = "Nombre del bucket S3 para uploads"
  value       = aws_s3_bucket.uploads.bucket
}

output "cloudwatch_log_group" {
  description = "Nombre del grupo de logs en CloudWatch"
  value       = aws_cloudwatch_log_group.api_logs.name
}
