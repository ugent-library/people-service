# build stage
FROM golang:1.21-alpine AS build

RUN go install github.com/jackc/tern/v2@latest
RUN cp $GOPATH/bin/tern /usr/local/bin/

# final stage
FROM alpine:latest

WORKDIR /migrations

COPY etc/migrations .
COPY docker/dbmigrate.sh .
COPY --from=build /usr/local/bin/tern /usr/local/bin/tern
RUN chmod +x /migrations/dbmigrate.sh

ENV PGDSN $PGDSN

CMD "/migrations/dbmigrate.sh"
