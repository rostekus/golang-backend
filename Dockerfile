FROM alpine:3.15 as root
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:1.17 as builder
WORKDIR /server
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/main.go

FROM scratch as final
COPY --from=root /etc/passwd /etc/passwd
COPY --from=root /etc/group /etc/group
WORKDIR /server
COPY --from=builder /server/.env .env
COPY --chown=1001:1001 --from=builder /server/main /main
USER app
EXPOSE 8080
ENTRYPOINT ["/main"]
