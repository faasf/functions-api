FROM alpine:3.17

WORKDIR /app

COPY build/ .

CMD [ "/app/functions-api" ]