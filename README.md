
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
eval $(minikube docker-env)
docker build -t lyraproj/nodejs .
# we could do this sort of thing
# kubectl run lyra-dep --image=lyraproj/nodejs --image-pull-policy=Never
#
# but actually we have namespace, pod and service files ready to go
kubectl apply -f local-docker-reg
# easiest way to test (should really be doing this from outside the cluster - and outside the service)
kubectl exec --namespace local-docker-registry-test -it hello-node-js -- /bin/bash
curl 127.0.0.1:8080
```