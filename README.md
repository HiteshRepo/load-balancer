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
kubectl create -f build/restSimpleApp/deploy.yaml
```

6. Get the pods
```bash
kubectl get po -l app=rest-simple-app -o wide

NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE                              NOMINATED NODE   READINESS GATES
rest-simple-app-6c9dc776f5-gwnvq   2/2     Running   0          59s   10.244.0.4   eksabm-testbed1-core-dev-node-0   <none>           <none>
rest-simple-app-6c9dc776f5-hkp7w   2/2     Running   0          59s   10.244.0.2   eksabm-testbed1-core-dev-node-0   <none>           <none>
rest-simple-app-6c9dc776f5-vfnbx   2/2     Running   0          59s   10.244.0.3   eksabm-testbed1-core-dev-node-0   <none>           <none>
```

7. Curl the pods using their IP and listening router port
```bash
curl 10.244.0.4:9090 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-gwnvq

curl 10.244.0.3:9090 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-vfnbx

curl 10.244.0.2:9090 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-hkp7w
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
kubectl create -f build/proxyApp/deploy.yaml

kubectl get po -o wide -l app=proxy-app

NAME                         READY   STATUS    RESTARTS   AGE    IP             NODE                              NOMINATED NODE   READINESS GATES
proxy-app-6c6bd559d7-rx7m2   2/2     Running   0          3m4s   10.244.0.254   eksabm-testbed1-core-dev-node-0   <none>           <none>
```

3. Deploy the proxy svc

```bash
kubectl create -f build/proxyApp/svc.yaml

kubectl get svc

NAME                       TYPE           CLUSTER-IP       EXTERNAL-IP                                   PORT(S)                               AGE
proxy-svc                  ClusterIP      10.108.215.193   <none>                                        9092/TCP                              3m52s
```

4. Request to add all pod-urls to the proxy
```bash
curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.4:9090"}' -w "\n"
proxy registered

curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.2:9090"}' -w "\n"
proxy registered

curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.3:9090"}' -w "\n"
proxy registered
```

5. Request to serve proxy
```bash
curl 10.108.215.193:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-gwnvq

curl 10.108.215.193:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-gwnvq
```

Note: does not matter how many times you hit the proxy, request is redirected to same po [the first one all the time].
Can you spot the issue in the code in the branch 'setup-proxy'?
What we have achieved though is to setup a proxy that redirects the request to the pod(s) behind it.
Now we would like to move on to achieve load balancing!!

6. Before we move on to load balancing, lets setup periodic health check of the urls we add behind our proxy

the proxy app's image 'hiteshpattanayak/proxy-app:2.0' has the health check feature setup.
Also the image has been updated in 'proxyApp-deploy.yaml'.

Deploy the latest proxy app
```bash
kubectl apply -f build/proxyApp/deploy.yaml

deployment.apps/proxy-app configured

kubectl get po -l app=proxy-app

NAME                        READY   STATUS    RESTARTS   AGE
proxy-app-6c6bd559d7-rx7m2   2/2     Running   0          29s
```

7. Repeat the 'addUrl' requests again

```bash
curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.4:9090"}' -w "\n"
proxy registered

curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.3:9090"}' -w "\n"
proxy registered

curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.2:9090"}' -w "\n"
proxy registered

## adding an invalid url
curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://1.2.3.4:9090"}' -w "\n"
proxy registered
```

8. check logs of the proxy app pod
```bash
kubectl logs -f proxy-app-555465b5b-xwz42 -c proxy-app

[GIN] 2023/04/16 - 12:33:08 | 200 |     420.363µs |       127.0.0.1 | POST     "/addurl"
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
[GIN] 2023/04/16 - 12:33:14 | 200 |    1.198732ms |       127.0.0.1 | POST     "/addurl"
16047-09+00 47:23:23 -> url(http://10.244.0.4:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.2:9090) is available.
[GIN] 2023/04/16 - 12:33:19 | 200 |     292.743µs |       127.0.0.1 | POST     "/addurl"
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
16047-09+00 47:23:23 -> url(http://1.2.3.4:9090) is unavailable.
[GIN] 2023/04/16 - 12:33:25 | 200 |     337.648µs |       127.0.0.1 | POST     "/addurl"
16047-09+00 47:23:23 -> url(http://10.244.0.4:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.2:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
16047-09+00 47:23:23 -> url(http://1.2.3.4:9090) is unavailable.
16047-09+00 47:23:23 -> url(http://10.244.0.4:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.2:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
16047-09+00 47:23:23 -> url(http://1.2.3.4:9090) is unavailable.
16047-09+00 47:23:23 -> url(http://10.244.0.4:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.2:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
16047-09+00 47:23:23 -> url(http://1.2.3.4:9090) is unavailable.
16047-09+00 47:23:23 -> url(http://10.244.0.4:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.2:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
16047-09+00 47:23:23 -> url(http://1.2.3.4:9090) is unavailable.
16047-09+00 47:23:23 -> url(http://10.244.0.4:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.2:9090) is available.
16047-09+00 47:23:23 -> url(http://10.244.0.3:9090) is available.
16047-09+00 47:23:23 -> url(http://1.2.3.4:9090) is unavailable.
```

9. Lets deploy `hiteshpattanayak/proxy-app:3.0` which contains implementation of `round robin` load balancing strategy.

10. Deploy the latest proxy app
```bash
kubectl apply -f build/proxyApp/deploy.yaml

deployment.apps/proxy-app configured
```

11. Repeat the 'addUrl' requests again, we do this because we do not have a persistent stirage as yet.
All the urls are stored in-memory.

```bash
curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.4:9090"}' -w "\n"
proxy registered

curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.3:9090"}' -w "\n"
proxy registered

curl 10.108.215.193:9092/addurl --header "Content-Type: application/json" --request "POST" --data '{"url": "http://10.244.0.2:9090"}' -w "\n"
proxy registered
```

12. Request to serve proxy now

```bash
curl 10.108.215.193:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-vfnbx

curl 10.108.215.193:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-hkp7w

curl 10.108.215.193:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-gwnvq

curl 10.108.215.193:9092 -w "\n"
responding from a simple rest app with pod name: rest-simple-app-6c9dc776f5-vfnbx
```

Note: The load balancer redirects requests across pods in a round-robin manner.

What if we could redirect requests in a weighted fashion?
Could set an url's (server's) capability to handle more requests/load and the load balancer could redirect request based on that?
Lets find out in the next steps.