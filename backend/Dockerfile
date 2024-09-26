FROM golang:1.22-alpine AS build

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

COPY . /app
WORKDIR /app
RUN apk add build-base 
RUN go build -trimpath -o server ./cmd/server
RUN go build -trimpath -o user-register ./cmd/user-register
RUN go build -trimpath -o sync-tool ./cmd/sync-tool

FROM golang:1.22-alpine

WORKDIR /app

COPY --from=build /app/server /usr/local/bin
COPY --from=build /app/user-register /usr/local/bin
COPY --from=build /app/sync-tool /usr/local/bin

CMD [ "server" ]
