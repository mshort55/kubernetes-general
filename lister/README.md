# lister
This application lists pods and deployments in kube-system namespace.

Run manually:
1. go run main.go --kubeconfig ~/.kube/config

Run in k8s:
1. kubectl create clusterrole listpodsdeploy --resource pods,deployments --verb list
2. kubectl create clusterrolebinding lister --clusterrole lister --serviceaccount default:default
3. go build
4. docker build -t lister:0.1.x .
5. docker tag lister:0.1.x docker.io/mshort55/lister:0.1.x
6. docker push docker.io/mshort55/lister:0.1.x
7. vim lister.yaml (update image tag with new version)
8. kubectl apply -f lister.yaml