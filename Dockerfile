FROM alpine:3.17

WORKDIR /app
RUN mkdir -p internal/config

COPY build/ .
COPY internal/config/config.yml ./internal/config/config.yml

CMD [ "/app/functions-api" ]