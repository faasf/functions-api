FROM alpine:3.17

WORKDIR /app

RUN mkdir config && touch config/config.yml

COPY build/ .

CMD [ "/app/functions-api" ]