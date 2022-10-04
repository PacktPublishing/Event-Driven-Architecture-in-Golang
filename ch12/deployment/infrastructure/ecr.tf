// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ecr_authorization_token
data aws_ecr_authorization_token token {}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecr_repository
resource aws_ecr_repository services {
  for_each     = toset(var.services)
  name         = each.key
  force_delete = true
}

// https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs/resources/registry_image
resource docker_registry_image services {
  for_each = toset(var.services)
  name     = "${aws_ecr_repository.services[each.key].repository_url}:latest"

  build {
    context    = "../.."
    dockerfile = "docker/Dockerfile.microservices"

    build_args = {
      service = each.key
    }
  }
}

output ecr_url {
  value = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.region}.amazonaws.com"
}
