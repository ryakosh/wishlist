FROM golang:1.14.4-alpine3.12

WORKDIR /api
COPY . .
RUN go mod download && go get github.com/cosmtrek/air
ENTRYPOINT [ "air" ]