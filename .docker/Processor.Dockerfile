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
RUN go build -o ./bin/processor cmd/processor/processor.go

FROM alpine as processor

COPY --from=builder /usr/local/src/bin/processor ./processor
COPY config.docker /config

ENTRYPOINT ["./processor"]

#ENTRYPOINT ["/bin/sh","-c","sleep infinity"]



