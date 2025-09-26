# docker buildx build -t kng1259/kmgm-consumer ./consumer
# docker buildx build -t kng1259/kmgm-producer ./producer
# docker push kng1259/kmgm-consumer
# docker push kng1259/kmgm-producer

docker buildx build -t poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main/kmgm-consumer:latest consumer/
docker buildx build -t poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main/kmgm-producer:latest producer/
docker push poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main/kmgm-consumer:latest
docker push poc4k-tsnode1b.ovng.dev.myovcloud.com/docker-main/kmgm-producer:latest