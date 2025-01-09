# https://registry.terraform.io/providers/hashicorp/aws/latest/docs#argument-reference

provider "aws" {
  region  = var.region
  profile = var.profile

  assume_role {
    role_arn     = "arn:aws:iam::846072081665:role/TerraformDeployDevRole"
    session_name = "TerraformSession"
  }

  default_tags {
    tags = {
      Product = var.product
      Env     = var.env
    }
  }
}
