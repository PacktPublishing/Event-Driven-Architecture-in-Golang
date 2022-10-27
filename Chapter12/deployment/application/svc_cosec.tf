// https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password
resource random_password cosec {
  length = 16
}

// https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource
// https://www.terraform.io/language/resources/provisioners/local-exec
resource null_resource init_cosec_db {
  provisioner "local-exec" {
    command     = "psql --file sql/init_service_db.psql -v db=$DB -v user=$USER -v pass=$PASS ${local.db_conn}/postgres"
    environment = {
      DB   = "cosec"
      USER = "cosec_user"
      PASS = random_password.cosec.result
    }
  }
  depends_on = [
    null_resource.init_db,
    random_password.cosec
  ]
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret_v1
resource kubernetes_secret_v1 cosec {
  metadata {
    name      = "cosec-secrets"
    namespace = local.project
  }

  data = {
    PG_CONN = "host=${local.db_host} port=${local.db_port} dbname=cosec user=cosec_user password=${random_password.cosec.result} search_path=cosec,public"
  }
  depends_on = [
    kubernetes_namespace_v1.namespace,
    null_resource.init_cosec_db
  ]
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment_v1
resource kubernetes_deployment_v1 cosec {
  metadata {
    name      = "cosec"
    namespace = local.project
    labels    = {
      app = "cosec"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        "app.kubernetes.io/name" = "cosec"
      }
    }
    template {
      metadata {
        name   = "cosec"
        labels = {
          "app.kubernetes.io/name" = "cosec"
        }
      }
      spec {
        hostname = "cosec"
        container {
          name              = "cosec"
          image             = "${local.aws_ecr_url}/cosec:latest"
          image_pull_policy = "Always"
          env_from {
            config_map_ref {
              name = "common-config-map"
            }
          }
          env_from {
            secret_ref {
              name = "cosec-secrets"
            }
          }
          port {
            protocol       = "TCP"
            container_port = 80
          }
          port {
            protocol       = "TCP"
            container_port = 9000
          }
          liveness_probe {
            http_get {
              path = "/liveness"
              port = 80
            }
            initial_delay_seconds = 3
            period_seconds        = 5
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_namespace_v1.namespace,
    kubernetes_config_map_v1.common,
    kubernetes_secret_v1.cosec,
    kubernetes_service_v1.nats
  ]
}
