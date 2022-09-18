// https://registry.terraform.io/modules/terraform-aws-modules/security-group/aws/4.13.0
module security_group {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.13.0"

  name   = "${var.project}-sg"
  vpc_id = module.vpc.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 1
      protocol    = "TCP"
      to_port     = 65365
      cidr_blocks = "${var.allowed_cidr_block},${var.vpc_cidr_block}"
    },
    {
      from_port   = -1
      protocol    = "icmp"
      to_port     = -1
      cidr_blocks = "${var.allowed_cidr_block},${var.vpc_cidr_block}"
    }
  ]

  egress_with_cidr_blocks = [
    {
      from_port   = 0
      protocol    = "-1"
      to_port     = 0
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  tags = {
    "kubernetes.io/cluster/${var.project}": "shared"
  }
}
