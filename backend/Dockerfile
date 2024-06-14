FROM golang:1.22-alpine AS build

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

COPY . /app
WORKDIR /app
RUN apk add build-base 
RUN go build -trimpath -o server ./cmd/server/main.go
RUN go build -trimpath -o user-register ./cmd/register/main.go

FROM golang:1.22-alpine

WORKDIR /app
COPY --from=build /app/server /app
COPY --from=build /app/user-register /app

CMD [ "./server" ]
