FROM golang:1.17 AS build

ARG GIT_DESC=undefined

WORKDIR /go/src/github.com/mysteriumnetwork/ethdiff
COPY . .
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-s -w -extldflags "-static" -X main.version='"$GIT_DESC" ./cmd/ethdiff
ADD https://curl.haxx.se/ca/cacert.pem /certs.crt
RUN chmod 0644 /certs.crt

FROM scratch AS arrange
COPY --from=build /go/src/github.com/mysteriumnetwork/ethdiff/ethdiff /
COPY --from=build /certs.crt /etc/ssl/certs/ca-certificates.crt

FROM scratch
COPY --from=arrange / /
USER 9999:9999
ENTRYPOINT ["/ethdiff"]
