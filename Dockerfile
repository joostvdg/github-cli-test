FROM golang:1.14 as builder
WORKDIR /go/src/github-cli-test
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ghcs

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV CONFIG_LOCATION "/config.toml"
ENV TOKEN_LOCATION "apitoken"
COPY --from=builder /go/src/github-cli-test/ghcs /ghcs
EXPOSE 1323
CMD ["/ghcs"]