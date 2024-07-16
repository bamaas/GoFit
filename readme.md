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
USERS='[{"email": "demo@gofit.nl", "password": "gofit123"}]'
docker run -e USERS=${USERS} --name gofit -v ${PWD}/data:/data -p 8080:8080 bamaas/gofit:latest
```

### Kubernetes

A Helm chart is available at [deploy/chart/gofit](deploy/chart/gofit)

