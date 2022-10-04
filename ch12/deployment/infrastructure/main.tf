variable region {
  description = "AWS Region"
  type        = string
}

variable services {
  description = "List of MallBots microservices"
  type        = list(string)
  default     = ["baskets", "cosec", "customers", "depot", "ordering", "notifications", "payments", "search", "stores"]
}

variable project {
  description = "Project name"
  type        = string
  default     = "mallbots"
}

variable allowed_cidr_block {
  description = "CIDR allowed to access public resources (application, bastion, ...) Example: \"Your Public IP\"/32"
  type        = string
}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity
data aws_caller_identity current {}

// https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string
resource random_string suffix {
  length  = 8
  special = false
}

output region {
  description = "AWS Region"
  value       = var.region
}

output project {
  description = "Project name"
  value       = var.project
}

output allowed_cidr_block {
  value = var.allowed_cidr_block
}

output services {
  value = var.services
}
