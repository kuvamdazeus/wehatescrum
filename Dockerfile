FROM golang:bookworm

WORKDIR /app
COPY . .

ARG gh_token
ENV TOKEN $gh_token
ENV PORT 8080

EXPOSE 8080

WORKDIR /app/tmp
RUN chmod +x ../entrypoint.sh
RUN ../entrypoint.sh

WORKDIR /app
ENTRYPOINT ["go", "run", "./src"]