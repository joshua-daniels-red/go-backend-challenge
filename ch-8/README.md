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

---

## 🔎 Observability

### Port-forward Prometheus:

```bash
kubectl port-forward svc/prometheus 9090:9090
```

→ [http://localhost:9090](http://localhost:9090)

---

### Port-forward Grafana:

```bash
kubectl port-forward svc/grafana 3000:3000
```

→ [http://localhost:3000](http://localhost:3000)
Login: `admin` / `admin`

You'll find a pre-provisioned dashboard under **"Observability Dashboard"** visualizing:

* Events produced to Redpanda
* Events consumed from Redpanda
* Events processed successfully / failed
* Stream input rates

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
├── setup.sh                     # Run-all setup script
```

---

## 📌 Notes

* Prometheus scrapes `/metrics` from both `producer` and `consumer`
* Cassandra stores domain and user stats in `goanalytics` keyspace
* Dashboard provisioning requires correct datasource names (`Prometheus`)
* All resources are provisioned using `kubectl apply`, no Helm

---

## ✅ Teardown

To remove everything:

```bash
kubectl delete -f k8s/
kubectl delete pvc --all
```

---

## ✅ Status

All services are deployed and observable within Kubernetes. This chapter completes the infrastructure transformation of the project.
