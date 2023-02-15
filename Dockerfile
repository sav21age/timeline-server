FROM golang:1.18.10-alpine3.17 AS builder
RUN go version
COPY . /github.com/sav21age/timeline-server/
WORKDIR /github.com/sav21age/timeline-server/
RUN go mod download
RUN GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /github.com/sav21age/timeline-server/.bin/app .
COPY --from=0 /github.com/sav21age/timeline-server/config config/
COPY --from=0 /github.com/sav21age/timeline-server/template template/

CMD ["./app"]