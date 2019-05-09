FROM golang:1.12

MAINTAINER  blademainer <blademainer@gmail.com>

ENV APP_HOME="/etcd-sync"

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

ADD ./conf/etcd-sync.yml $APP_HOME/conf/etcd-sync.yml
ADD ./conf/logger.yml $APP_HOME/conf/logger.yml

ADD ./bin/etcd-sync* $APP_HOME


ENTRYPOINT ["./etcd-sync"]

