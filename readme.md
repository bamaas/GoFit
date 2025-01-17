# GoFit

A bodyweight tracking application.

## Demo

Demo @ [https://gofit.basmaas.nl](https://gofit.basmaas.nl).

Username: `demo@gofit.nl`

Password: `gofit123`

It's running on the free tier of Azure WebApps, so expect slow performance.

## How to run

### Docker

```bash
GID=$(id -g)
UID=$(id -u)
USERS='[{"email": "demo@gofit.nl", "password": "gofit123"}, {"email": "user@gofit.nl", "password": "gofit123"}]'
DATA_DIR="./data"
mkdir ${DATA_DIR}

docker run \
--rm \
-d \
--user ${UID}:${GID} \
-e USERS="${USERS}" \
--name gofit \
-v ${DATA_DIR}:/data \
-p 8080:8080 \
bamaas/gofit:latest

docker logs -f gofit
```

### Kubernetes

A Helm chart is available at [deploy/chart/gofit](deploy/chart/gofit)

## TODO

* Fix issue when initializing app with 1 user.
