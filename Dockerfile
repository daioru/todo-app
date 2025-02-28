# Builder
FROM golang:1.24.0-alpine3.21 AS builder

ARG GITHUB_PATH=github.com/daioru/todo-app
WORKDIR /home/${GITHUB_PATH}

COPY . .

RUN apk add --no-cache make
RUN make build
RUN make migrate_build

# Server
FROM alpine:latest AS server

LABEL org.opencontainers.image.source=https://${GITHUB_PATH}

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ARG GITHUB_PATH=github.com/daioru/todo-app

COPY --from=builder /home/${GITHUB_PATH}/bin/todo-app .
COPY --from=builder /home/${GITHUB_PATH}/bin/migration .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations
# Копируем entrypoint.sh
COPY --from=builder /home/${GITHUB_PATH}/entrypoint.sh .

# Копирование .env
COPY .env .env

# Делаем скрипт исполняемым
RUN chmod +x /root/entrypoint.sh  

EXPOSE 8080

# Запускаем скрипт при старте контейнера
ENTRYPOINT ["/root/entrypoint.sh"]  
