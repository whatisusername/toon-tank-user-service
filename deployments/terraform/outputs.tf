################################################################################
# Lambda Function
################################################################################

output "lambda_function_arn" {
  description = "The ARN of the Lambda Function"
  value       = module.lambda_with_image.lambda_function_arn
}

output "lambda_function_invoke_arn" {
  description = "The Invoke ARN of the Lambda Function"
  value       = module.lambda_with_image.lambda_function_invoke_arn
}

output "lambda_function_name" {
  description = "The name of the Lambda Function"
  value       = module.lambda_with_image.lambda_function_name
}

output "lambda_function_qualified_arn" {
  description = "The ARN identifying your Lambda Function Version"
  value       = module.lambda_with_image.lambda_function_qualified_arn
}

output "lambda_function_qualified_invoke_arn" {
  description = "The Invoke ARN identifying your Lambda Function Version"
  value       = module.lambda_with_image.lambda_function_qualified_invoke_arn
}

output "lambda_alias_name" {
  description = "The name of the Lambda Function Alias"
  value       = module.lambda_alias_release.lambda_alias_name
}

output "lambda_alias_arn" {
  description = "The ARN of the Lambda Function Alias"
  value       = module.lambda_alias_release.lambda_alias_arn
}

output "lambda_alias_invoke_arn" {
  description = "The ARN to be used for invoking Lambda Function from API Gateway"
  value       = module.lambda_alias_release.lambda_alias_invoke_arn
}

################################################################################
# CloudWatch Logs
################################################################################

output "lambda_cloudwatch_log_group_arn" {
  description = "The ARN of the Cloudwatch Log Group"
  value       = module.lambda_with_image.lambda_cloudwatch_log_group_arn
}

output "lambda_cloudwatch_log_group_name" {
  description = "The name of the Cloudwatch Log Group"
  value       = module.lambda_with_image.lambda_cloudwatch_log_group_name
}
