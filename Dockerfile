FROM golang:1.17-alpine3.15 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o jira-telegram .

FROM alpine:3.15

WORKDIR /bin

COPY --from=builder --chown=1001:1001 /build/jira-telegram .

USER 1001

ENTRYPOINT [ "/bin/jira-telegram" ]