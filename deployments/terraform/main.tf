# Data Source: aws_iam_role
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_role

data "aws_iam_role" "user_auth" {
  name = "UserAuthHandlerLambdaRole"
}

# Data Source: aws_ecr_image
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ecr_image

data "aws_ecr_image" "main" {
  repository_name = format("%s/%s", var.env, lower(var.product))
  image_digest    = var.image_digest == null ? null : var.image_digest
  image_tag       = var.image_digest == null ? var.image_tag : null
  most_recent     = var.image_digest == null && var.image_tag == null ? true : null
}

# module: lambda
# https://registry.terraform.io/modules/terraform-aws-modules/lambda/aws/latest

module "lambda_with_image" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.17.0"

  function_name  = format("%s-%s", var.name, var.env)
  description    = var.description
  create_role    = false
  lambda_role    = data.aws_iam_role.user_auth.arn
  memory_size    = var.memory_size
  publish        = true
  timeout        = var.timeout
  image_uri      = data.aws_ecr_image.main.image_uri
  create_package = false
  package_type   = "Image"
  architectures  = ["x86_64"]

  environment_variables = {
    SECRET_NAME = format("%s-cognito-secrets-%s", lower(var.product), var.env)
  }

  use_existing_cloudwatch_log_group = false
  cloudwatch_logs_retention_in_days = 30
  cloudwatch_logs_skip_destroy      = false
  cloudwatch_logs_log_group_class   = "STANDARD"

  timeouts = {
    update = "3m"
  }
}

# submodule: alias
# https://registry.terraform.io/modules/terraform-aws-modules/lambda/aws/latest/submodules/alias

module "lambda_alias_release" {
  source  = "terraform-aws-modules/lambda/aws//modules/alias"
  version = "~> 7.17.0"

  name             = var.alias_name
  function_name    = module.lambda_with_image.lambda_function_name
  function_version = module.lambda_with_image.lambda_function_version

  depends_on = [
    module.lambda_with_image
  ]
}
