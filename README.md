# Marzban API Gateway

## Description

This API Gateway provides gRPC interface to complete limited set of interactions with Marzban API.

## Deployment

2 options for deployment are available:

1. Use docker image [luckycomet55/marzban-api-gtw](https://hub.docker.com/repository/docker/luckycomet55/marzban-api-gtw)
2. Build from source

Both methods require a set of environmental variables to be provided:
```
ADMIN_USERNAME    -- username for admin in a Marzban instance
ADMIN_PASSWORD    -- password for admin in a Marzban instance
BASE_URL          -- Marzban instance base url
ENV               -- dev/prod. dev uses structured text logging to stdout, prod uses JSON
PORT              -- port for gRPC server to listen to. Defaults to 8343
```


### 1. Using docker image

After setting all the environment variables run

```
docker pull luckycomet55/marzban-api-gtw:<TAG>
docker run --env-file <ENV_FILE_PATH> -p <OUT_PORT>:<PORT>
```

Assert that <PORT> is the same as provided to container envs.


### 2. Building from source

You must have go >=1.24 preinstalled on your system.
Clone the repository, then run
```
make setup
make build BUILD_DIR=<BUILD_DIR> BINARY_NAME=<BINARY_NAME>
```

You may not provide variables, they default to
```
BUILD_DIR=./build
BINARY_NAME=marzban-api-gtw
```
