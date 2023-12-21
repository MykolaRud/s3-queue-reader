FROM golang:1.20-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash

COPY ["go.mod", "go.sum", "./"]
COPY ["cmd/", "./cmd"]
COPY ["handlers", "./handlers"]
COPY ["infrastructures", "./infrastructures"]
COPY ["interfaces", "./interfaces"]
COPY ["models", "./models"]
COPY ["repositories", "./repositories"]

RUN go mod download
RUN go build -o ./bin/watcher cmd/watcher/watcher.go

FROM alpine as watcher

COPY --from=builder /usr/local/src/bin/watcher ./watcher
COPY config.docker /config

#CMD [ "./watcher" ]
ENTRYPOINT [ "./watcher" ]
#ENTRYPOINT ["/bin/sh","-c","sleep infinity"]



