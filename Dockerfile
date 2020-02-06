FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY *.go go.mod go.sum ./

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pidtemp . 
RUN ls -l

FROM alpine

COPY --from=builder /app/pidtemp /go/bin/pidtemp

CMD ["/go/bin/pidtemp"] 

# This dir needs to be mounted and contain a config.toml file.
WORKDIR /app/workdir

