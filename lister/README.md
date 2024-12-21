# lister
This lists out cluster objects.

Run manually:
1. go run main.go --kubeconfig ~/.kube/config

Run in k8s:
1. go build
2. docker build -t lister:0.1.x .
3. docker tag lister:0.1.x docker.io/mshort55/lister:0.1.x
4. update image tag in lister.yaml
5. kubectl apply -f lister.yaml