# Chapter 8: Kubernetes Deployment

This chapter containerizes and deploys all components of the analytics system into a Kubernetes environment using Minikube. It marks the final step in transitioning the application from local Docker Compose to production-grade Kubernetes infrastructure.

---

## 🧱 Components Deployed

* ✅ **Cassandra** (StatefulSet + Service)
* ✅ **Redpanda** (StatefulSet + Service)
* ✅ **Producer** (Deployment)
* ✅ **Consumer** (Deployment)
* ✅ **Prometheus** (for metrics scraping)
* ✅ **Grafana** (for observability dashboards)
* ✅ **Schema Initializers**:
  * Cassandra keyspace + tables
  * Redpanda topic (`wikipedia.protobuf`)

---

## ✨ Quickstart

Make sure you have:

* [Minikube](https://minikube.sigs.k8s.io/docs/start/) installed
* Docker image built inside Minikube (see below)

### 1. Start Minikube

```bash
minikube start --cpus=4 --memory=6g
```

### 2. Use Minikube's Docker Daemon

```bash
eval $(minikube docker-env)
```

### 3. Build Images in Minikube

```bash
docker build -t producer-app:latest -f Dockerfile.producer .
docker build -t consumer-app:latest -f Dockerfile.consumer .
```

### 4. Run the Full Setup

```bash
chmod +x setup.sh
./setup.sh
```

This script ensures proper initialization order:

* Cassandra starts before the consumer
* Redpanda is ready before the topic init and producer
* Dashboards and data sources are provisioned for Grafana

After this script completes, **all pods and services will be running**, including:

* Working `/stats` endpoint on the Consumer
* Writing events to Cassandra
* Populated Grafana dashboards

---

## 🔎 Accessing Services

Since this setup is local, services are exposed using **`NodePort`** and can be reached using the Minikube IP. You can run:

```bash
minikube service list
```

Or directly open in browser:

```bash
minikube service grafana
minikube service prometheus
minikube service consumer
```

You can also retrieve the IP manually:

```bash
minikube ip
```

Then access services via:

* Grafana:     `http://<minikube-ip>:30300`
* Prometheus:  `http://<minikube-ip>:30900`
* Consumer:    `http://<minikube-ip>:30080/stats`

---

## 📊 Observability Dashboard

Grafana is provisioned with a full dashboard under:

> **"Observability Dashboard"**

It visualizes:

* Events produced to Redpanda
* Events consumed from Redpanda
* Events processed successfully / failed
* Stream input rates

Log in with:

* **Username**: `admin`
* **Password**: `admin`

---

## 🤖 CI/CD Integration

Chapter 8 includes a dedicated CI pipeline using **GitHub Actions** and a **KinD (Kubernetes-in-Docker)** cluster. This allows for full validation of:

* Docker builds of the producer and consumer apps
* Cluster setup via `setup.sh`
* Integration tests hitting `/metrics` and `/stats`
* Full teardown using `teardown.sh`

The pipeline ensures your Kubernetes configuration and deployments remain healthy and verifiable on every push.

---

## 📦 Project Structure (Kubernetes)

```
ch-8/
├── k8s/
│   ├── cassandra/                # Cassandra StatefulSet + Service
│   ├── cassandra-init/          # Job to create keyspace + tables
│   ├── redpanda/                # Redpanda StatefulSet + Service
│   ├── redpanda-init/           # Job to create topic
│   ├── config/                  # Shared environment ConfigMap
│   ├── producer/                # Producer Deployment + Service
│   ├── consumer/                # Consumer Deployment + Service
│   ├── prometheus/              # Prometheus config + deployment
│   └── grafana/                 # Grafana dashboards + data sources
├── db/init.cql                  # Cassandra schema file
├── monitoring/                  # Dashboard + Prometheus config
├── setup.sh                     # Run-all setup script
├── teardown.sh                  # Wipes all Kubernetes resources
├── .github/workflows/ci-ch8.yml# CI pipeline definition
```

---

## 📌 Teardown

To clean up everything:

```bash
chmod +x teardown.sh
./teardown.sh
```

This removes:

* All Kubernetes resources (`kubectl delete -f k8s/`)
* Any persistent volumes (`kubectl delete pvc --all`)
* Associated ConfigMaps

---

## ✅ Status

All services are deployed, connected, and observable within Kubernetes. This chapter completes the infrastructure transformation of the project, and the entire stack can now be managed through `setup.sh`, tested through CI with KinD, and accessed via `minikube service` commands for full local orchestration.
