terraform {
  backend "s3" {
    dynamodb_table = "deployment-terraform-lock"
    region         = "eu-west-2"
    bucket         = "deployment-terraform-123456789"
    key            = "terraform.tfstate"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.11.0"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "eu-central-1"
}

provider "aws" {
  alias  = "infra-account"
  region = "eu-central-1"
  assume_role {
    role_arn = "arn:aws:iam::0987654321:role/assumable_role"
  }
}
locals {
  app_name    = "bruno-beans"
  app_version = 1
  app_float   = 1.45

  cba_base_domain  = var.cba_base_domain
  tasks            = jsondecode("[{\"name\":\"datetime\",\"image\":\"datetime-image-path\",\"env\":[{\"name\":\"ALWAYS_LOG_WARNINGS_STDERR\",\"value\":\"1\"},{\"name\":\"SERVER_ROLE_MAY_RUN_EXPERIMENTS\",\"value\":\"0\"}]}]")
  rollout_strategy = jsondecode("{\"Steps\":[{\"Name\":\"staging\",\"Traffic\":{\"Old\":100,\"New\":0}},{\"Name\":\"vanguard\",\"Traffic\":{\"Old\":90,\"New\":10}},{\"Name\":\"full on\",\"Traffic\":{\"Old\":0,\"New\":100}}]}")
  iam_policy_names = ["rds_policy_1", "rds_policy_2"]
}

locals {
  tags = {
    Terraform = "true"
    App       = local.app_name
  }
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}
module "bruno-beans-7132aaa" {
  source = "git::https://github.com/module"

  app_name         = local.app_name
  app_version      = local.cell_version
  cba_base_domain  = local.cba_base_domain
  tasks            = local.tasks
  iam_policy_names = local.iam_policy_names
  cba_environment  = local.environment
  tag = {
    deployment_sha1 = "7132aaa4-6db3-4cae-8d36-8e903fd06698"
  }
}
