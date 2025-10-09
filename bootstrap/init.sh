# kubectl apply -f applications/
# sleep 5
kubectl apply -f sigstore/
sleep 60
kubectl apply -f cilium/ccnp.yaml
kubectl apply -f cert-manager/
kubectl apply -f sigstore/
kubectl apply -f prometheus/