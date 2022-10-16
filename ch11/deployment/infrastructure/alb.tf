variable lb_image_repository {
  description = "AWS Load Balancer image host See: https://docs.aws.amazon.com/eks/latest/userguide/add-ons-images.html"
  type        = string
  default     = "602401143452.dkr.ecr.us-east-1.amazonaws.com"
}

variable lb_service_account_name {
  type    = string
  default = "aws-load-balancer-controller"
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_account_v1
resource kubernetes_service_account_v1 lb_service_account {
  metadata {
    name      = var.lb_service_account_name
    namespace = "kube-system"
    labels    = {
      "app.kubernetes.io/name"      = "aws-load-balancer-controller"
      "app.kubernetes.io/component" = "controller"
    }
    annotations = {
      "eks.amazonaws.com/role-arn" = module.vpc_cni_irsa.iam_role_arn
    }
  }
}

// https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release
resource helm_release lb {
  name = "load-balancer"

  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"

  set {
    name  = "clusterName"
    value = var.project
  }

  set {
    name  = "serviceAccount.create"
    value = "false"
  }

  set {
    name  = "serviceAccount.name"
    value = var.lb_service_account_name
  }

  set {
    name  = "region"
    value = var.region
  }

  // See: https://docs.aws.amazon.com/eks/latest/userguide/add-ons-images.html
  set {
    name  = "image.repository"
    value = "${var.lb_image_repository}/amazon/aws-load-balancer-controller"
  }

  set {
    name  = "vpcId"
    value = module.vpc.vpc_id
  }

  depends_on = [
    module.eks,
    kubernetes_service_account_v1.lb_service_account
  ]
}
