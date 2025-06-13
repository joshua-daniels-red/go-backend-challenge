#!/bin/bash
set -e

echo "🧾 Creating required ConfigMaps..."

kubectl create configmap cassandra-schema \
  --from-file=init.cql=db/init.cql --dry-run=client -o yaml | kubectl apply -f -

kubectl create configmap prometheus-config \
  --from-file=prometheus.yml=monitoring/prometheus.yml \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create configmap grafana-dashboard-config \
  --from-file=monitoring/provisioning/dashboards/definitions/observability.json \
  --dry-run=client -o yaml | kubectl apply -f -


kubectl create configmap grafana-dashboard-definitions \
  --from-file=monitoring/provisioning/dashboards/definitions/observability.json \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create configmap grafana-datasource \
  --from-file=monitoring/provisioning/datasources/datasource.yml \
  --dry-run=client -o yaml | kubectl apply -f -

echo "🚀 Applying core services..."
kubectl apply -f k8s/cassandra/
kubectl wait --for=condition=Ready pod/cassandra-0 --timeout=90s


echo "🧰 Initializing Cassandra schema..."
kubectl apply -f k8s/cassandra-init/
kubectl wait --for=condition=complete job/cassandra-init --timeout=180s|| {
  echo "❌ Cassandra init job failed or timed out"
  exit 1
}

# 🔍 Actively check that the keyspace exists
echo "⏳ Probing Cassandra keyspace readiness..."
until kubectl exec cassandra-0 -- cqlsh -e "DESCRIBE KEYSPACE goanalytics" > /dev/null 2>&1; do
  echo "[WAIT] Keyspace not ready yet..."; sleep 2;
done


echo "🪞 Deploying Redpanda..."
kubectl apply -f k8s/redpanda/
kubectl wait --for=condition=Ready pod/redpanda-0 --timeout=90s


echo "🎯 Creating Redpanda topic..."
kubectl apply -f k8s/redpanda-init/
kubectl wait --for=condition=complete job/redpanda-topic-init --timeout=60s


echo "🪖 Deploying producer and consumer..."
kubectl apply -f k8s/config/
kubectl apply -f k8s/producer/
kubectl apply -f k8s/consumer/


echo "📊 Deploying Prometheus and Grafana..."
kubectl apply -f k8s/prometheus/
kubectl apply -f k8s/grafana/


echo "🌟 Setup complete. You can now port-forward to Grafana and view your dashboard:"
echo "kubectl port-forward svc/grafana 3000:3000"
