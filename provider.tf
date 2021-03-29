terraform {
  required_providers {
    aws-1 = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}
provider "aws" {
  region = "ap-northeast-1"
}
