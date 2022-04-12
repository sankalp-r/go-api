FROM golang:1.17-alpine3.13 as builder
RUN mkdir /app
WORKDIR /app
RUN apk add --update --no-cache ca-certificates git
## Add this go mod download command to pull in any dependencies
COPY pkg ./pkg
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/app/main ./api

WORKDIR /app/api

RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags "-s" -a -installsuffix cgo .

FROM alpine:3.13.5


ENV APP_USER app
ENV APP_HOME /app
ENV GROUP_ID 3000
ENV USER_ID 1000

RUN addgroup -g ${GROUP_ID} -S ${APP_USER} && adduser -u ${USER_ID} -S ${APP_USER} -G ${APP_USER}
RUN mkdir -p $APP_HOME


WORKDIR /app
COPY --from=builder /app/api/main ./

RUN chown -R $APP_USER:$APP_USER $APP_HOME

USER $APP_USER

ENTRYPOINT ["./main"]