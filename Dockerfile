FROM golang:1.22-alpine

RUN apk update
RUN apk upgrade
RUN apk add git bash

WORKDIR /app
COPY . .

ARG gh_token
ENV TOKEN $gh_token
ENV PORT 8000

EXPOSE 8000

WORKDIR /app/tmp
RUN chmod +x ../entrypoint.sh
RUN ../entrypoint.sh

WORKDIR /app
RUN go build -o bin/main src/*

# =========== STAGE 2 (RUN BUILD) ===========
FROM alpine

WORKDIR /app
COPY --from=0 /app/tmp tmp
COPY --from=0 /app/bin bin

ENTRYPOINT ["./bin/main"]