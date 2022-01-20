FROM alpine:latest

ADD ./scripts/docker-entrypoint.sh /docker-entrypoint.sh
ADD ./scripts/publish.sh /publish.sh

RUN apk update &&\
  apk add --update curl python3 &&\
  curl https://raw.githubusercontent.com/rabbitmq/rabbitmq-server/v3.9.11/deps/rabbitmq_management/bin/rabbitmqadmin -o /usr/bin/rabbitmqadmin &&\
  chmod +x /usr/bin/rabbitmqadmin &&\
  chmod +x /docker-entrypoint.sh &&\
  chmod +x /publish.sh

# possible environment variables with defaults
ENV RABBIT_HOST=127.0.0.1 \
  RABBIT_PORT=15672 \
  RABBIT_USER=guest \
  RABBIT_PASSWORD=guest \
  RABBIT_VHOST=/ \
  RABBIT_EXCHANGE=amq.default \
  RABBIT_QUEUE=nap-tasks \
  RABBIT_PAYLOAD="{}"

ENTRYPOINT ["/docker-entrypoint.sh"]
