#!/bin/bash
set -e

echo "🧨 Tearing down Kubernetes resources..."

# Delete all resources applied from k8s manifests
find k8s/ -name "*.yaml" -o -name "*.yml" | xargs -n1 -I{} kubectl delete -f {} --ignore-not-found
##kubectl delete -f k8s/ --ignore-not-found

# Delete ConfigMaps used by Grafana, Prometheus, Cassandra
echo "🧹 Deleting ConfigMaps..."
kubectl delete configmap cassandra-schema \
  grafana-dashboard-definitions \
  grafana-dashboard-config \
  grafana-datasource \
  prometheus-config --ignore-not-found

# Delete persistent volume claims
echo "🗑️ Deleting PVCs..."
kubectl delete pvc --all

echo "✅ Teardown complete. Your Minikube cluster is clean."
