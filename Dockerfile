FROM alpine:edge

RUN apk add --no-cache go

EXPOSE 8080

WORKDIR /app

ENV SRC_DIR=/app
RUN mkdir -p $SRC_DIR

ADD . $SRC_DIR

RUN apk add --no-cache tzdata
ENV TZ America/New_York

RUN cd $SRC_DIR; go get ./...; go build -o /app/re2playground

ENTRYPOINT ["/app/re2playground"]
