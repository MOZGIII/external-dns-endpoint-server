FROM golang:1.16.5 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-s -extldflags "-static"' -mod=vendor -o build .

FROM scratch
COPY --from=builder /app/build /usr/local/bin/app
CMD ["/usr/local/bin/app"]
