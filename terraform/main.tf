terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "2.13.1"
    }
    minikube = {
      source  = "scott-the-programmer/minikube"
      version = "0.3.10"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.29.0"
    }
  }
}

provider "minikube" {
  kubernetes_version = "v1.28.3"
}

provider "helm" {
  kubernetes {
    config_path    = "~/.kube/config"
    config_context = "terraform-provider-minikube-acc-docker"
  }
}

provider "kubernetes" {
  host = minikube_cluster.docker.host

  client_certificate     = minikube_cluster.docker.client_certificate
  client_key             = minikube_cluster.docker.client_key
  cluster_ca_certificate = minikube_cluster.docker.cluster_ca_certificate
}
