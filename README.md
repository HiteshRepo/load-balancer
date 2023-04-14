## Steps for the setup

1. Install `task` cli tool:
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

2. cd into root directory:
```bash 
cd ../load-balancer
```

3. compile, package and publish restSimpleApp:
```bash
task compile-restSimpleApp

task package-restSimpleApp

task publish-restSimpleApp
```

Note: skip steps #1, #2 and #3 since, the app's image has already been pushed to dockerhub.

4. Setup a kubernetes cluster, easier way to do this is via minikube: https://kubernetes.io/docs/tutorials/hello-minikube/

5. Deploy the app
```bash
kubectl create -f build/restSimpleApp-deploy.yaml
```

6. Get the pods
```bash
kubectl get po -o wide

rest-simple-app-964f5d5bd-q4m99             2/2     Running   0          20m     10.244.0.32    minikube   <none>           <none>
rest-simple-app-964f5d5bd-q9jqm             2/2     Running   0          20m     10.244.0.30    minikube   <none>           <none>
rest-simple-app-964f5d5bd-qmdrt             2/2     Running   0          20m     10.244.0.31    minikube   <none>           <none>
```

7. Curl the pods using their IP and listening router port
```bash
curl 10.244.0.32:9090 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-964f5d5bd-q4m99

curl 10.244.0.31:9090 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-964f5d5bd-qmdrt

curl 10.244.0.30:9090 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-964f5d5bd-q9jqm
```