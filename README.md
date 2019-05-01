
# Namespace secrets-sandbox

Locking down a single secret to a single service account

```sh 
# Clean up before starting (adapted from https://stackoverflow.com/a/50396865)
kubectl delete daemonsets,replicasets,services,deployments,pods,rc,secrets,serviceaccounts,namespaces,roles,rolebindings  --all
# wait for that to complete (when error, this loop will exit)
while kubectl get namespaces secrets-sandbox && sleep 2; do :; done
# Begin:
# Create a namespace, service account, secret, role (and bind role to service account) and pod in one shot (apply does this is alpha order of filenames)
# secrets must be base-64 encoded e.g. "echo -n 'joebloggs' | base64" -> am9lYmxvZ2dz
kubectl apply -f secrets-sandbox/
# sh on to the container for a look around
kubectl exec --namespace secrets-sandbox -it my-pod -- /bin/bash
# we can see this password (or the username or )
cat /etc/mark-secrets/password
# and available we have made them available as environment variables on that container (see 6-pod.yaml)
echo $SECRET_PASSWORD
# the azure-sa can get pods cos we gave that role but not, for example, get services
kubectl --namespace secrets-sandbox --as azure-sa auth can-i get pods
# we have locked down access to a single secret (see all/5-role.yaml) from that service account
kubectl --namespace secrets-sandbox --as azure-sa auth can-i get secret/azure-sa-secret # Yes!
kubectl --namespace secrets-sandbox --as azure-sa auth can-i get secret/other-secret # "no - no RBAC policy matched" \o/
```

# Running a local container in minikube

Using the local docker registry (on the host machine) 

```sh
minikube start
minikube dashboard
# create namespace and deploy
kubectl create namespace testnamespace
kubectl get namespaces
kubectl --namespace testnamespace create deployment nginx --image=nginx 
kubectl get deployments
kubectl --namespace testnamespace get deployments
kubectl delete namespace testnamespace
# docker - two interdependent services
eval "$(docker-machine env -u)" # undo eval minikube
docker build -t lyraproj/customer customer-service && docker build -t lyraproj/address address-service
docker run -t -p 8081:8081 lyraproj/address
curl http://localhost:8081 
docker run -t -p 8080:8080 --env ADDRESS_URI=http://host.docker.internal:8081 lyraproj/customer
curl http://localhost:8080 
# in k8s
eval $(minikube docker-env)
docker build -t lyraproj/customer customer-service && docker build -t lyraproj/address address-service
kubectl apply -f local-docker-reg
# export customer_pod_name=$(kubectl --namespace local-docker-registry-test get pods | grep customer-deployment | head -1 | awk {'print $1'})
# kubectl exec --namespace local-docker-registry-test -it $customer_pod_name /bin/bash
# curl 127.0.0.1:8080
curl http://$(minikube ip):30003
curl http://$(minikube ip):30002
# get logs
kubectl --namespace local-docker-registry-test get pods | grep address-deployment | awk '{print $1}' | head -1 | xargs kubectl --namespace local-docker-registry-test logs
kubectl --namespace local-docker-registry-test get pods | grep customer-deployment | awk '{print $1}' | head -1 | xargs kubectl --namespace local-docker-registry-test logs
```

prep

```sh 
minikube delete
minikube start
eval "$(docker-machine env -u)"
docker build -t lyraproj/customer customer-service && docker build -t lyraproj/address address-service
eval "$(minikube docker-env)"
docker build -t lyraproj/customer customer-service && docker build -t lyraproj/address address-service
```

# Helm

```sh
brew install kubernetes-helm
# init helm and install tiller
helm init --history-max 200
helm repo update              # Make sure we get the latest list of charts
helm install stable/mysql
```
