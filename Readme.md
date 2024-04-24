# health

This repository contains health check service. 

It polls provided endpoints and checks if they are healthy.

There is a list of response processors that can be used to check the response of the endpoint:
 * LogProcessor - logs the response with appropriate log level
 * EmailProcessor - sends an email if the response is not healthy (mocked)

We can add more processors in the future.

## Table of Contents
1. [How to run](#how-to-run)
    1. [Prerequisites](#prerequisites)
    2. [docker-compose](#docker-compose)
    3. [GoLand](#goland)
    4. [Go](#go)
    5. [Minikube (IaC)](#minikube-iac)
2. [Supported flags](#supported-flags)
3. [Sources discovery](#sources-discovery)
4. [Future improvements](#future-improvements)
5. [Mock services](#mock-services)

## How to run

---

#### Prerequisites
  * [docker](https://docs.docker.com/get-docker/) running
  * [docker-compose](https://docs.docker.com/compose/install/) installed
  * [go](https://golang.org/doc/install) installed, version 1.22 or higher recommended
  * [helm](https://helm.sh/docs/intro/install/) installed
  * [terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) installed

### docker-compose
The easiest way to run the service is to use docker-compose.

```bash
make compose_up
```

This will:
  * run the tests
  * run helm template
  * run the service
  * run 3 additional mock services (with wiremock) that the health check service will poll

### GoLand

There is a run configuration saved in the [.run](.run) folder. You can use it in GoLand IDE and run the service from there.
(Edit Configurations->Go Build)

### Go

You can run the runner using the following command:

```bash
go run cmd/main.go runner --sources-store-dir=.local/sources --sources-store-file=sources.json
```

If you combine it with docker-compose command your local runner will poll the mock services running in docker-compose.

### Minikube (IaC)

You can run the service in minikube using the following command:

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

This will:
  * create a minikube cluster (docker)
  * deploy the service
  * deploy 3 additional mock services (with wiremock) that the health check service will poll

You can inspect the cluster using tools like [k9s](https://github.com/derailed/k9s) or 
[kubectl](https://kubernetes.io/docs/reference/kubectl/overview/).

## Supported flags

---

Easiest way to see all supported flags is to run the service with the `--help` flag.

```bash
❯ go run cmd/main.go --help
NAME:o run cmd/main.go --help                                                                                                                                                                                                                                                                                                                               ─╯
   health-checker - run health-checker commands

USAGE:
   health-checker [global options] command [command options] 

COMMANDS:
   runner   run health checker
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

```bash
❯ go run cmd/main.go runner -h
NAME:o run cmd/main.go runner -h                                                                                                                                                                                                                                                                                                                            ─╯
   health-checker runner - run health checker

USAGE:
   health-checker runner [command options] [arguments...]

OPTIONS:
   --load FILE                           load configuration from yaml FILE [$LOAD]
   --default-health-check-path value     default health check path (default: "/health") [$DEFAULT_HEALTH_CHECK_PATH]
   --workers value                       number of workers (default: 1) [$WORKERS]
   --check-interval value                check interval (default: 1m0s) [$CHECK_INTERVAL]
   --log-level value                     log level (default: "info") [$LOG_LEVEL]
   --log-format value                    log format (default: "json") [$LOG_FORMAT]
   --sources-store-dir value             sources store directory (default: "/sources") [$SOURCES_STORE_DIR]
   --sources-store-file value            sources store file (default: "sources.yaml") [$SOURCES_STORE_FILE]
   --transient-errors-max-retries value  transient errors max retries (default: 3) [$TRANSIENT_ERRORS_MAX_RETRIES]
   --transient-errors-retry-wait value   transient errors retry wait (default: 1s) [$TRANSIENT_ERRORS_RETRY_WAIT]
   --http-client-timeout value           http client timeout (default: 5s) [$HTTP_CLIENT_TIMEOUT]
   --acceptable-response-time value      acceptable response time (default: 3s) [$ACCEPTABLE_RESPONSE_TIME]
   --help, -h                            show help
```

## Sources discovery

---

At the moment the sources are stored in a file. The file is in the yaml format and for local development it is stored 
in the `.local/sources` directory. When running the service in docker-compose the `sources-docker-compose.yaml` file 
is mounted as a volume. For local development you can use the `sources.yaml` file from the `.local/sources` directory.

Helm chart uses the `sources.yaml` file from the `deploy/charts/health/files` directory. It generates a configmap 
which is then mounted as a volume in the pod.

Loading the sources from a file is just one way of providing the sources. We can add more ways in the future:
 * load from external service 
 * load from kubernetes discovery (list services with appropriate labels attached)

**Reloading sources**

Runner is prepared to reload the sources from the injected `Loader` (config file loader, k8s labels scraper etc). 
It should be wired to allow the runner to reload the sources on specific interval (without the need to restart the service.

## Future improvements

---

There are a few things that can be improved in the future:
 * add more response processors (like slack, pagerduty etc)
 * add more sources discovery methods (like k8s labels scraper)
 * make the sources discovery more dynamic (reload sources on specific interval)
 * maybe introduce service as a sources provisioner (service that provides sources to the health check service)
   so that the health check service can have HPA and scale up/down based on needs. Provisioner would have to keep
   track of which sources should be provided to the health service (some FIFO queue).
 * add more tests
 * expose health check results via API and/or prometheus metrics
 * add tracing options (Grafana builds service map)

## Mock services

---

Mock services are provided using [Wiremock](http://wiremock.org/). 

There are 3 mock services that the health check service:
 * service-a - always returns 200 and all components are healthy
 * service-b - always returns 502
 * service-c - always returns 200 but one component is unhealthy
