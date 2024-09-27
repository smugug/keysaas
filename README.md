


# Build local image
./build-local-deploy-artifacts.sh
# Run local executable
./test  

# Ingress controller (node port: http:30010 https:30011)
<!-- helm repo add haproxytech https://haproxytech.github.io/helm-charts
helm repo update -->


# Create self-signed wildcards cert
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certificate/kubernetes.local.key.pem -out certificate/kubernetes.local.crt.pem \
  -subj "/CN=*.kubernetes.local/O=YourOrg" \
  -addext "subjectAltName=DNS:*.kubernetes.local,DNS:kubernetes.local"
  

# Create keysaas CRD and RBAC
kubectl apply -f artifacts/keysaas-crd-deployment.yaml
# Create namespace for deployed Keysaases
kubectl create namespace customer2
# Create database
kubectl apply -f artifacts/postgresql.yaml

kubectl create secret tls -n customer2 kubernetes-tls --cert=certificate/tls.crt --key=certificate/tls.key

helm install haproxy-kubernetes-ingress haproxytech/kubernetes-ingress \
  --create-namespace \
  --namespace haproxy-controller

# Deploy keysaas
kubectl apply -f artifacts/keysaastest.yaml

# TODO  
- keycloak
- rbac
- web gui
- monitor
