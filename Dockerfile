FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pidtemp . 
RUN ls -l

FROM scratch

COPY --from=builder /app/pidtemp /go/bin/pidtemp

WORKDIR /app/config
CMD ["/go/bin/pidtemp"] 

