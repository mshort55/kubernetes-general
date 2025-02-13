# expose-controller 

This is an example controller. It watches all cluster deployments and creates service and ingress
resources if they do not already exist. It will expose each deployment pod on port 80.

This controller will panic if a deployment is deleted. I plan to fix this in the future so that it deleted
ingresses and services upon deployment deletion.

Prerequisites:
You need to install an ingress controller in order for the ingress objects to work.
If using minikube, you can do so with this command:
1. minikube addons enable ingress

Procedure:
1. go build
2. ./expose-operator --kubeconfig ~/.kube/config
3. kubectl create deployment nginx --image nginx

You should see new ingress and service resources created with the same name as the newly created deployment.

Testing procedure:
1. kubectl get ingress -n <new-deployment-namespace>
2. note IP in address column
3. curl <address>:80/<deployment-name>