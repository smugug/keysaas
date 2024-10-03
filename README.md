## Install
### Add repo traefik and prometheus
helm repo add traefik https://traefik.github.io/charts
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

### Keysaas CRD and RBAC
kubectl apply -f artifacts/keysaas-crd-deployment.yaml
### Nnamespace for deployed Keysaases
kubectl create namespace customer2
### Database
kubectl apply -f artifacts/postgresql.yaml
### Prometheus
curl -sL https://github.com/prometheus-operator/prometheus-operator/releases/download/v0.69.0/bundle.yaml | kubectl create -f -
kubectl apply -f artifacts/prometheus/prometheus.yaml
### Metric server
kubectl apply -f artifacts/metrics-server/components.yaml --kubelet-insecure-tls
### Traefik (ingress controller)
helm install traefik traefik/traefik -f artifacts/traefik/values.yaml
### Cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml
kubectl apply -f artifacts/cert-manager/clusterissuer.yaml
### Web backend
kubectl apply -f artifacts/rbac-web/keysaasrole.yaml

## Oof
### Test keysaas
kubectl apply -f artifacts/keysaastest.yaml
### Build local image
./build-local-deploy-artifacts.sh
### Run local executable
./test.sh

## TODO

## NOTE
welp somehow you can't use axios to connect to other pods with cacert, request works well