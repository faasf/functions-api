FROM alpine:3.17

WORKDIR /app
RUN mkdir -p internal/config && touch internal/config/config.yml

COPY build/ .

CMD [ "/app/functions-api" ]