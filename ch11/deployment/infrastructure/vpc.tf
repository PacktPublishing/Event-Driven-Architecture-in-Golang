variable vpc_cidr_block {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable vpc_public_subnets {
  description = "List of public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable vpc_private_subnets {
  description = "List of private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable vpc_database_subnets {
  description = "List of database subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.5.0/24", "10.0.6.0/24"]
}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/availability_zones
data aws_availability_zones available {
  state = "available"
}

// https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/3.14.3
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.14.0"

  name = "${var.project}-vpc"

  cidr = var.vpc_cidr_block
  azs  = slice(data.aws_availability_zones.available.names, 0, 2)

  private_subnets  = var.vpc_private_subnets
  public_subnets   = var.vpc_public_subnets
  database_subnets = var.vpc_database_subnets

  enable_nat_gateway = true
  single_nat_gateway = true

  enable_dns_support   = true
  enable_dns_hostnames = true

  // Allows public access to the database
  create_database_subnet_group           = true
  create_database_subnet_route_table     = true
  create_database_internet_gateway_route = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.project}" = "shared"
    "kubernetes.io/role/elb"               = 1
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${var.project}" = "shared"
    "kubernetes.io/role/internal-elb"      = 1
  }
}

output vpc_cidr_block {
  value = var.vpc_cidr_block
}
