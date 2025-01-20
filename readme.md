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
ghcr.io/bamaas/gofit:0.0.1

docker logs -f gofit
```

### Kubernetes

A Helm chart is available.

```bash
helm repo add gofit https://bamaas.github.io/GoFit/
helm repo update
helm install gofit gofit/gofit
```

### Development setup

Start the devcontainer from within VSCode and you are good to go.

## TODO

* Fix issue when initializing app with 1 user.
