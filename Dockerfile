FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . /src/

RUN apk add --no-cache make bash gcc musl-dev
RUN make setup && make build


FROM alpine:3.22.1
COPY --from=build /src/build/marzban-api-gtw /app/marzban-api-gtw
CMD ["/app/marzban-api-gtw"]
