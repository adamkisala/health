FROM golang:1.22 as build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOSUMDB=off

WORKDIR /src/
COPY . /src/

RUN go build -mod=vendor -a -tags netgo --installsuffix netgo -ldflags '-w' -o health cmd/main.go


FROM alpine:3.19.1
RUN apk --no-cache add ca-certificates

RUN adduser -D -s /bin/sh app

COPY --from=build /src/health /bin/health
RUN chmod a+x /bin/health

USER app
