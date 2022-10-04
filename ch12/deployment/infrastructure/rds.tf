variable db_instance_type {
  description = "RDS serverless instance type"
  type        = string
  default     = "db.serverless"
}

variable db_family {
  description = "RDS serverless family"
  type        = string
  default     = "aurora-postgresql13"
}

variable db_username {
  description = "User name for the RDS PostgreSQL database"
  type        = string
}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/rds_engine_version
data aws_rds_engine_version postgres {
  engine = "aurora-postgresql"
}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_parameter_group
resource aws_db_parameter_group postgres {
  name   = "${var.project}-parameter-group"
  family = var.db_family
}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster_parameter_group
resource aws_rds_cluster_parameter_group postgres {
  name   = "${var.project}-cluster-parameter-group"
  family = var.db_family
}

// https://registry.terraform.io/modules/terraform-aws-modules/rds-aurora/aws/7.5.1
module "db" {
  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "~> 7.5.0"

  name                   = "${var.project}-db-cluster"
  instance_class         = var.db_instance_type
  engine                 = data.aws_rds_engine_version.postgres.engine
  engine_mode            = "provisioned"
  engine_version         = data.aws_rds_engine_version.postgres.version
  master_username        = var.db_username
  create_random_password = true
  port                   = 5432

  db_parameter_group_name         = aws_db_parameter_group.postgres.id
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.postgres.id

  instances = {
    primary = {}
  }

  serverlessv2_scaling_configuration = {
    min_capacity = 1
    max_capacity = 5
  }

  vpc_id  = module.vpc.vpc_id
  subnets = module.vpc.database_subnets

  allowed_cidr_blocks    = module.vpc.private_subnets_cidr_blocks
  db_subnet_group_name   = module.vpc.database_subnet_group_name
  create_db_subnet_group = false
  create_security_group  = false
  vpc_security_group_ids = [module.security_group.security_group_id]

  apply_immediately   = true
  skip_final_snapshot = true

  # This should never be set for a real production database!
  publicly_accessible = true
}

output db_endpoint {
  value = module.db.cluster_endpoint
}

output db_port {
  value = module.db.cluster_port
}

output db_username {
  value     = module.db.cluster_master_username
  sensitive = true
}

output db_password {
  value     = module.db.cluster_master_password
  sensitive = true
}

output db_conn {
  value     = "postgres://${module.db.cluster_master_username}:${module.db.cluster_master_password}@${module.db.cluster_endpoint}:${module.db.cluster_port}"
  sensitive = true
}
