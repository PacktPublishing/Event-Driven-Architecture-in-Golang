// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment_v1
resource kubernetes_deployment_v1 nats {
  metadata {
    name      = "nats"
    namespace = local.project
    labels    = {
      app = "nats"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        "app.kubernetes.io/name" = "nats"
      }
    }
    template {
      metadata {
        name   = "nats"
        labels = {
          "app.kubernetes.io/name" = "nats"
        }
      }
      spec {
        hostname       = "nats"
        restart_policy = "Always"
        container {
          image = "nats:2-alpine"
          name  = "nats"
          args  = ["-m", "8222", "-js", "-sd", "/var/lib/nats/data"]
          port {
            protocol       = "TCP"
            container_port = 4222
          }
          volume_mount {
            mount_path = "/var/lib/nats/data"
            name       = "jsdata"
          }
          liveness_probe {
            http_get {
              path = "/"
              port = 8222
            }
            initial_delay_seconds = 3
            period_seconds        = 5
          }
        }
        volume {
          name = "jsdata"
          persistent_volume_claim {
            claim_name = "jsdata"
          }
        }
      }
    }
  }
  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_v1
resource kubernetes_service_v1 nats {
  metadata {
    name      = "nats"
    namespace = local.project
    labels    = {
      app = "nats"
    }
  }
  spec {
    selector = {
      "app.kubernetes.io/name" = "nats"
    }
    port {
      protocol    = "TCP"
      port        = 4222
      target_port = 4222
    }
    type = "ClusterIP"
  }
  depends_on = [
    kubernetes_namespace_v1.namespace,
    kubernetes_deployment_v1.nats
  ]
}

// https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/persistent_volume_claim_v1
resource kubernetes_persistent_volume_claim_v1 nats {
  metadata {
    name      = "jsdata"
    namespace = local.project
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage : "100Mi"
      }
    }
  }
  depends_on = [
    kubernetes_namespace_v1.namespace
  ]
}
