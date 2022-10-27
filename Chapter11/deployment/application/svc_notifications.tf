// https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password
resource random_password notifications {
  length = 16
}

// https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource
// https://www.terraform.io/language/resources/provisioners/local-exec
resource null_resource init_notifications_db {
  provisioner "local-exec" {
    command     = "psql --file sql/init_service_db.psql -v db=$DB -v user=$USER -v pass=$PASS ${local.db_conn}/postgres"
    environment = {
      DB   = "notifications"
      USER = "notifications_user"
      PASS = random_password.notifications.result
    }
  }
  depends_on = [
    null_resource.init_db,
    random_password.notifications
  ]
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret_v1
resource kubernetes_secret_v1 notifications {
  metadata {
    name      = "notifications-secrets"
    namespace = local.project
  }

  data = {
    PG_CONN = "host=${local.db_host} port=${local.db_port} dbname=notifications user=notifications_user password=${random_password.notifications.result} search_path=notifications,public"
  }
  depends_on = [
    kubernetes_namespace_v1.namespace,
    null_resource.init_notifications_db
  ]
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment_v1
resource kubernetes_deployment_v1 notifications {
  metadata {
    name      = "notifications"
    namespace = local.project
    labels    = {
      app = "notifications"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        "app.kubernetes.io/name" = "notifications"
      }
    }
    template {
      metadata {
        name   = "notifications"
        labels = {
          "app.kubernetes.io/name" = "notifications"
        }
      }
      spec {
        hostname = "notifications"
        container {
          name              = "notifications"
          image             = "${local.aws_ecr_url}/notifications:latest"
          image_pull_policy = "Always"
          env_from {
            config_map_ref {
              name = "common-config-map"
            }
          }
          env_from {
            secret_ref {
              name = "notifications-secrets"
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
