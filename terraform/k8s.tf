resource "kubernetes_namespace" "monitoring" {
  metadata {
    name = "monitoring"
  }
}

resource "helm_release" "health" {
  name            = "health"
  chart           = "../deploy/charts/health"
  version         = "0.1.0"
  namespace       = "monitoring"
  wait            = true
  timeout         = 120
  reset_values    = true
  atomic          = true
  force_update    = true
  cleanup_on_fail = true
  values = [
    file("../deploy/charts/health/values.yaml")
  ]

  depends_on = [kubernetes_namespace.monitoring]

  set {
    name  = "runner.workers"
    value = "3"
  }
  set {
    name  = "image.tag"
    value = "0.1.3"
  }
}
