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
