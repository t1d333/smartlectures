FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN GOOS=linux go build -tags=jsoniter -o server ./cmd/easytex/main.go 

FROM debian:bookworm AS build-release-stage

ARG USERNAME=nonroot
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && chown -R $USER_UID:$USER_GID /home/$USERNAME

WORKDIR /

COPY --from=build-stage /app/server /nonroot/server

EXPOSE 8000

USER nonroot

ENTRYPOINT ["./nonroot/server"]
