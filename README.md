## Steps for the setup

### Simple REST app

The below app just gives out its own pod name in response on a GET call to its '/' endpoint.

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


### Proxy app

This app is designed to act as a proxy for any number of urls.

1. compile, package and publish proxyApp:
```bash
task compile-proxyApp

task package-proxyApp

task publish-proxyApp
```

Note: skip steps #1 since, the app's image has already been pushed to dockerhub.

2. Deploy the proxy app

```bash
kubectl create -f proxyApp-deploy.yaml

kubectl get po -o wide -l app=proxy-app

NAME                        READY   STATUS    RESTARTS   AGE   IP             NODE                              NOMINATED NODE   READINESS GATES
proxy-app-784f88697-lnfj8   2/2     Running   0          22m   10.244.0.187   eksabm-testbed1-core-dev-node-0   <none>           <none>
```

3. Deploy the proxy svc

```bash
kubectl create -f proxyApp-svc.yaml

kubectl get svc

NAME                       TYPE           CLUSTER-IP       EXTERNAL-IP                                   PORT(S)                               AGE
proxy-svc                  ClusterIP      10.110.72.7      <none>                                        9092/TCP                              22m
```

4. Request to add all pod-urls to the proxy
```bash
curl 10.110.72.7:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.32:9090"}' -w "\n"
proxy registered

curl 10.110.72.7:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.31:9090"}' -w "\n"
proxy registered

curl 10.110.72.7:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.30:9090"}' -w "\n"
proxy registered
```

5. Request to serve proxy
```bash
curl 10.110.72.7:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-964f5d5bd-q4m99

curl 10.110.72.7:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-964f5d5bd-q4m99
```

Note: does not matter how many times you hit the proxy, request is redirected to same po [the first one all the time].
Can you spot the issue in the code in the branch 'setup-proxy'?
What we have achieved though is to setup a proxy that redirects the request to the pod(s) behind it.
Now we would like to move on to achieve load balancing!!