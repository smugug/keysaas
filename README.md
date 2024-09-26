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
# Deploy keysaas
kubectl apply -f artifacts/keysaastest.yaml