terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # https://developer.hashicorp.com/terraform/language/backend/s3
  backend "s3" {
    region         = "us-west-2"
    bucket         = "toontank-terraform-state-us-west-2"
    encrypt        = true
    dynamodb_table = "toontank-terraform-locks"

    profile        = "TerraformBackend"
    assume_role = {
      role_arn     = "arn:aws:iam::846072081665:role/TerraformStateAccessRole"
      session_name = "TerraformSession"
    }
  }
}
