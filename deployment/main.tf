terraform {
  backend "s3" {
    key     = "common/lambda/cloudwatch-log-destination"
    encrypt = true
  }
}

locals {
  common_tags = {
    Owner       = "global"
    Environment = "production"
  }
}

data "terraform_remote_state" "master" {
  backend = "s3"
  config {
    bucket   = "terraform.bhavik.io"
    key      = "common/master"
    region   = "${var.aws_default_region}"
    profile  = "${var.profile}"
    role_arn = "arn:aws:iam::${var.operations_account_id}:role/${var.role_name}"
  }
}

provider "aws" {
  alias   = "master"
  region  = "${var.aws_default_region}"
  version = "~> 2.7.0"
  profile = "${var.profile}"
}

provider "aws" {
  alias   = "identity"
  region  = "${var.aws_default_region}"
  version = "~> 2.7.0"
  profile = "${var.profile}"

  assume_role {
    role_arn     = "arn:aws:iam::${data.terraform_remote_state.master.identity_account_id}:role/${var.role_name}"
    session_name = "terraform"
  }
}

provider "aws" {
  alias   = "operations"
  region  = "${var.aws_default_region}"
  version = "~> 2.7.0"
  profile = "${var.profile}"

  assume_role {
    role_arn     = "arn:aws:iam::${data.terraform_remote_state.master.operations_account_id}:role/OrganizationAccountAccessRole"
    session_name = "terraform"
  }
}

provider "aws" {
  alias   = "development"
  region  = "${var.aws_default_region}"
  version = "~> 2.7.0"
  profile = "${var.profile}"

  assume_role {
    role_arn     = "arn:aws:iam::${data.terraform_remote_state.master.development_account_id}:role/${var.role_name}"
    session_name = "terraform"
  }
}

provider "aws" {
  alias   = "production"
  region  = "${var.aws_default_region}"
  version = "~> 2.7.0"
  profile = "${var.profile}"

  assume_role {
    role_arn     = "arn:aws:iam::${data.terraform_remote_state.master.production_account_id}:role/${var.role_name}"
    session_name = "terraform"
  }
}

module "lambda_master" {
  source              = "./modules/lambda"
  lambda_version      = "${var.lambda_version}"
  log_destination_arn = "${data.terraform_remote_state.master.log_destination_arn}"
  kms_key_arn         = "${data.terraform_remote_state.master.default_kms_key_arn}"
  tags                = "${merge(local.common_tags, var.tags)}"

  providers = {
    aws = "aws.master"
  }
}

module "lambda_identity" {
  source              = "./modules/lambda"
  lambda_version      = "${var.lambda_version}"
  log_destination_arn = "${data.terraform_remote_state.master.log_destination_arn}"
  kms_key_arn         = "${data.terraform_remote_state.master.default_kms_key_arn}"
  tags                = "${merge(local.common_tags, var.tags)}"

  providers = {
    aws = "aws.identity"
  }
}

module "lambda_operations" {
  source              = "./modules/lambda"
  lambda_version      = "${var.lambda_version}"
  log_destination_arn = "${data.terraform_remote_state.master.log_destination_arn}"
  kms_key_arn         = "${data.terraform_remote_state.master.default_kms_key_arn}"
  tags                = "${merge(local.common_tags, var.tags)}"

  providers = {
    aws = "aws.operations"
  }
}

module "lambda_development" {
  source              = "./modules/lambda"
  lambda_version      = "${var.lambda_version}"
  log_destination_arn = "${data.terraform_remote_state.master.log_destination_arn}"
  kms_key_arn         = "${data.terraform_remote_state.master.default_kms_key_arn}"
  tags                = "${merge(local.common_tags, var.tags)}"

  providers = {
    aws = "aws.development"
  }
}

module "lambda_production" {
  source              = "./modules/lambda"
  lambda_version      = "${var.lambda_version}"
  log_destination_arn = "${data.terraform_remote_state.master.log_destination_arn}"
  kms_key_arn         = "${data.terraform_remote_state.master.default_kms_key_arn}"
  tags                = "${merge(local.common_tags, var.tags)}"

  providers = {
    aws = "aws.production"
  }
}