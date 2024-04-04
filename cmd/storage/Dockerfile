FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -o server ./cmd/storage/main.go 

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/server /server

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["./server"]
