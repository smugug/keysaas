# Build local image
./build-local-deploy-artifacts.sh
# Run local executable
./test  

# Ingress controller (node port: http:30010 https:30011)
helm repo add traefik https://traefik.github.io/charts
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update


# Create keysaas CRD and RBAC
kubectl apply -f artifacts/keysaas-crd-deployment.yaml
# Create namespace for deployed Keysaases
kubectl create namespace customer2
# Create database
kubectl apply -f artifacts/postgresql.yaml
# Deploy prometheus
curl -sL https://github.com/prometheus-operator/prometheus-operator/releases/download/v0.69.0/bundle.yaml | kubectl create -f -
kubectl apply -f artifacts/prometheus/prometheus.yaml
# Deploy traefik
helm install traefik traefik/traefik -f artifacts/traefik/values.yaml
# Deploy cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml
kubectl apply -f artifacts/cert-manager/clusterissuer.yaml

kubectl apply -f artifacts/rbac-web/keysaasrole.yaml

# Deploy keysaas
kubectl apply -f artifacts/keysaastest.yaml

# TODO
- web gui
- monitor

# NOTE
welp somehow you can't use axios to connect to other pods with cacert, request works well
