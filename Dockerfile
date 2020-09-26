FROM golang:1.15

ARG EMAIL
ENV EMAIL=$EMAIL

ARG PASSWORD
ENV PASSWORD=$PASSWORD

COPY . /Mailing

WORKDIR /Mailing

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /Mailing/cmd/app

RUN go build .

CMD go run .