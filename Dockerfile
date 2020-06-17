FROM golang:1.14.2 as builder

MAINTAINER Someone in the Build_n_Release Team.

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o exporter gitlabgoexporter.go

FROM scratch

USER 65534:65534
COPY --from=builder /app/exporter /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/exporter"]