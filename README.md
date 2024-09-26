# Build local image
./build-local-deploy-artifacts.sh
# Run local executable
./test  

# Create keysaas CRD and RBAC
kubectl apply -f artifacts/keysaas-crd-deployment.yaml
# Create namespace for deployed Keysaases
kubectl create namespace customer2
# Create database
kubectl apply -f artifacts/mysql.yaml
# Wait for database
kubectl get pods -n customer2
# Ingress controller (node port: http:30010 https:30011)
helm repo add haproxytech https://haproxytech.github.io/helm-charts
helm repo update
helm install haproxy-kubernetes-ingress haproxytech/kubernetes-ingress \
  --create-namespace \
  --namespace haproxy-controller

# Deploy keysaas
kubectl apply -f artifacts/keysaastest.yaml