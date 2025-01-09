variable "region" {
  description = "Provision in which AWS region."
  type        = string
  default     = "us-west-2"
}

variable "product" {
  description = "The production name. It's for the tag."
  type        = string
  default     = "ToonTank"
}

variable "profile" {
  type        = string
  description = "The AWS profile to use for the deployment."
  default     = null
}

variable "env" {
  description = "The environment name. Should be dev/stag/prod"
  type        = string
  default     = null

  validation {
    condition     = can(regex("^(dev|stag|prod)$", var.env))
    error_message = "The env should be dev, stag or prod."
  }
}

################################################################################
# Lambda Function
################################################################################

variable "name" {
  type        = string
  description = "Unique name for your Lambda Function."
  default     = "userAuthHandler"
}

variable "description" {
  type        = string
  description = "Description of what your Lambda Function does."
  default     = "Handle user sign up/in. Interact with Cognito."
}

variable "memory_size" {
  type        = number
  description = "Amount of memory in MB your Lambda Function can use at runtime."
  default     = 128
}

variable "timeout" {
  type        = number
  description = "Amount of time your Lambda Function has to run in seconds."
  default     = 3
}

variable "alias_name" {
  type        = string
  description = "Name for the alias you are creating."
  default     = "release"
}

################################################################################
# ECR Image
################################################################################

variable "image_digest" {
  type        = string
  description = "Sha256 digest of the image manifest."
  default     = null
}

variable "image_tag" {
  type        = string
  description = "Tag associated with this image."
  default     = null
}
