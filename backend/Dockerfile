FROM golang:1.22.1-alpine as builder

WORKDIR /src/
COPY go.mod /src/
RUN go mod download
COPY go.sum /src/

COPY *.go /src/
COPY discord/*.go /src/discord/
COPY games/*.go /src/games/
RUN CGO_ENABLED=0 go build -o /build/myapp

FROM alpine:latest

COPY --from=builder /build/myapp /pb/myapp

# uncomment to copy the local migrations dir into the image
COPY ./migrations /pb/migrations

# uncomment to copy the local pb_hooks dir into the image
COPY ./pb_hooks /pb/pb_hooks

EXPOSE 8080

# start PocketBase
CMD ["/pb/myapp", "serve", "--http=0.0.0.0:8080"]