# Docker and kubernetes
This area contains all necessary files to build FOXDEN docker images and
deployment files for Kubernetes ifrastricture.

### Docker builds
```
# build docker image for Frontend service
# please note last argument should be path to this area
docker build --build-arg srv=Frontend -t ghcr.io/chesscomputing/frontend .

# build docker image for Frontend service for a specific tag
docker build --build-arg srv=Frontend --build-arg tag=v0.0.1 -t ghcr.io/chesscomputing/frontend .
```

### Kubernetes setup
```
# create foxden namespace
kubectl create ns foxden

# create foxden secret
kubectl create secret generic foxden-secrets --from-file=.foxden.yaml --dry-run=client -o yaml | kubectl apply --namespace=foxden -f -

# create foxden schemas secret
kubectl create secret generic schema-secrets --from-file=schemas/test.json --from-file=schemas/ID1A3.json --from-file=schemas/ID3A.json --from-file=schemas/ID4B.json --dry-run=client -o yaml | kubectl apply --namespace=foxden -f -

# deploy services
kubectl apply -f /Authz.yaml
kubectl apply -f DataBookkeeping.yaml
kubectl apply -f DataDiscovery.yaml
kubectl apply -f DataManagement.yaml
kubectl apply -f Frontend.yaml
kubectl apply -f MLHub.yaml
kubectl apply -f PublicationService.yaml
kubectl apply -f SpecScansService.yaml
```
