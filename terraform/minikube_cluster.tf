resource "minikube_cluster" "docker" {
  driver       = "docker"
  cluster_name = "terraform-provider-minikube-acc-docker"
  addons = [
    "default-storageclass",
    "storage-provisioner"
  ]
}
