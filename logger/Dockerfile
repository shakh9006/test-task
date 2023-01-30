FROM golang:alpine

ENV APP_NAME logger

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

ADD go.mod .
ADD go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o ${APP_NAME} ./cmd/logger/main.go

CMD ./${APP_NAME}