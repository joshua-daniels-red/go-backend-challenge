# Chapter 6: Observability & Protobuf

In this module, we enhance our system with two major capabilities:

* **Protobuf** for structured, efficient message serialization
* **Observability** using Prometheus and Grafana to monitor system behavior

---

## 🔄 Features Implemented

### ✅ Protobuf Integration

* Introduced `.proto` schemas for stream event messages
* Generated and used Protobuf types to serialize and deserialize messages
* Added schema files to the Docker image to enable Redpanda Console decoding
* Created a new Redpanda topic with optimized settings for Protobuf messages

### ✅ Observability with Prometheus & Grafana

* Integrated Prometheus metrics endpoint (`/metrics`) in both producer and consumer services
* Deployed Prometheus and Grafana via Docker Compose
* Added custom application metrics:

  * `events_consumed_from_stream_total`
  * `events_produced_to_redpanda_total`
  * `events_consumed_from_redpanda_total`
  * `events_processed_successfully_total`
  * `events_failed_to_process_total`
* Metrics exposed via Prometheus and confirmed accessible via `/metrics`
* Created Grafana dashboards manually to visualize key metrics (e.g., using `rate(...)` queries)

---

## 🧪 CI/CD Pipeline Enhancements

* Used GitHub Actions to test, vet, lint, and format Go code
* Ensured test coverage for new metrics logic
* Docker image build and push included for both producer and consumer
* Integration with Prometheus/Grafana tested locally using Docker

---

## 🛠️ Run Locally

```bash
docker compose down -v
docker compose up --build
```

* Access the app at:

  * **Consumer app:** [http://localhost:8080](http://localhost:8080)
  * **Prometheus:** [http://localhost:9090](http://localhost:9090)
  * **Grafana:** [http://localhost:3000](http://localhost:3000) — Login with `admin/admin`

---

## 🧩 Manual Dashboard Queries (Grafana)

Example queries used in Grafana panels:

```promql
rate(events_consumed_from_stream_total[$__rate_interval])
rate(events_produced_to_redpanda_total[$__rate_interval])
rate(events_consumed_from_redpanda_total[$__rate_interval])
rate(events_processed_successfully_total[$__rate_interval])
rate(events_failed_to_process_total[$__rate_interval])
```

These help visualize ingestion, processing, and failure rates in near real-time.

---

## ✅ Completion Summary

* [x] Protobuf schemas and serialization in place
* [x] Redpanda integration with Protobuf messages
* [x] Prometheus metrics tracked and exposed
* [x] Grafana dashboards built manually for metric visualization
* [x] CI pipeline fully functional with test, format, and Docker build steps

---