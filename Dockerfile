FROM golang:1.17-alpine as builder
WORKDIR /go/src/gitlab.wal.hds.com/CSE/chillfs/kubernetes/shared-services/reserve-version/
COPY go.mod go.sum ./ 
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /go/bin/reserve-version 

FROM alpine:3.13.6
WORKDIR /app
ENV PORT="8080"
COPY --from=builder /go/bin/reserve-version .
ENTRYPOINT [ "/app/reserve-version"]