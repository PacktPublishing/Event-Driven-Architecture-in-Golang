// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity
data aws_caller_identity current {}

locals {
  services             = data.terraform_remote_state.infra.outputs.services
  vpc_cidr_block       = data.terraform_remote_state.infra.outputs.vpc_cidr_block
  allowed_cidr_block   = data.terraform_remote_state.infra.outputs.allowed_cidr_block
  region               = data.terraform_remote_state.infra.outputs.region
  aws_ecr_url          = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.terraform_remote_state.infra.outputs.region}.amazonaws.com"
  eks_cluster_id       = data.terraform_remote_state.infra.outputs.eks_cluster_id
  eks_vpc_cni_role_arn = data.terraform_remote_state.infra.outputs.eks_vpc_cni_role_arn
  project              = data.terraform_remote_state.infra.outputs.project
  db_conn              = data.terraform_remote_state.infra.outputs.db_conn
  db_host              = data.terraform_remote_state.infra.outputs.db_endpoint
  db_port              = data.terraform_remote_state.infra.outputs.db_port
}
