REGISTRY=poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main
TAG=0.1.0

# build and push images
docker buildx build -t $REGISTRY/kmgm-consumer:$TAG consumer/
docker buildx build -t $REGISTRY/kmgm-producer:$TAG producer/
docker push $REGISTRY/kmgm-consumer:$TAG 
docker push $REGISTRY/kmgm-producer:$TAG

# generate SBOMs
trivy image --image-src remote,docker -f cyclonedx $REGISTRY/kmgm-consumer:$TAG -o kmgm-consumer.cdx.json
trivy image --image-src remote,docker -f cyclonedx $REGISTRY/kmgm-producer:$TAG -o kmgm-producer.cdx.json

# sign images
cosign sign -y $REGISTRY/kmgm-consumer:$TAG
cosign sign -y $REGISTRY/kmgm-producer:$TAG

# attest and push SBOM to registry
cosign attest -y --key keys/cosign.key --type cyclonedx --predicate kmgm-consumer.cdx.json $REGISTRY/kmgm-consumer:$TAG
cosign attest -y --key keys/cosign.key --type cyclonedx --predicate kmgm-consumer.cdx.json $REGISTRY/kmgm-producer:$TAG

# verify images from registry
cosign verify \
  $REGISTRY/kmgm-consumer:$TAG \
  $REGISTRY/kmgm-producer:$TAG \
  --certificate-identity=kngondajob@gmail.com --certificate-oidc-issuer=https://login.microsoftonline.com | jq .

# verify attestation and scan sbom
cosign verify-attestation --key keys/cosign.pub --type cyclonedx $REGISTRY/kmgm-consumer:$TAG > consumer.cdx.intoto.jsonl
  # --certificate-identity=kngondajob@gmail.com --certificate-oidc-issuer=https://login.microsoftonline.com > consumer.cdx.intoto.jsonl
cosign verify-attestation --key keys/cosign.pub --type cyclonedx $REGISTRY/kmgm-producer:$TAG > producer.cdx.intoto.jsonl
  # --certificate-identity=kngondajob@gmail.com --certificate-oidc-issuer=https://login.microsoftonline.com > producer.cdx.intoto.jsonl
trivy sbom consumer.cdx.intoto.jsonl
trivy sbom producer.cdx.intoto.jsonl

# verify image from bitnamisecure/mongodb
cosign verify --key https://app-catalog.vmware.com/.well-known/cosign.pub bitnamisecure/mongodb:latest --insecure-ignore-tlog | jq .

# then pull images after verifying

# # verify images after cosign save
# cosign save poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main/kmgm-consumer:latest --dir ./image/
# cosign verify --local-image image/ --certificate-identity=kngondajob@gmail.com --certificate-oidc-issuer=https://login.microsoftonline.com | jq .

# sign and verify helm chart from oci registry
trivy fs --scanners misconfig,secret helm-chart/ --output helmscan.txt
helm package helm-chart/
helm push kmgm-app-$TAG.tgz oci://$REGISTRY
cosign sign -y $REGISTRY/kmgm-app:$TAG --key keys/cosign.key
cosign verify $REGISTRY/kmgm-app:$TAG --key keys/cosign.pub | jq .

# helm pull oci://poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main/kmgm-app:0.1.0
