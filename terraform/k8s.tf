resource "kubernetes_namespace" "monitoring" {
  metadata {
    name = "monitoring"
  }
}

resource "helm_release" "weaviate_health" {
  name            = "weaviate-health"
  chart           = "../deploy/charts/weaviate-health"
  version         = "0.1.0"
  namespace       = "monitoring"
  wait            = true
  timeout         = 120
  reset_values    = true
  atomic          = true
  force_update    = true
  cleanup_on_fail = true
  values = [
    file("../deploy/charts/weaviate-health/values.yaml")
  ]

  depends_on = [kubernetes_namespace.monitoring]
}
