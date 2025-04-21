# docker-rabbitmqadmin

Simple go program that could've been a bash script to publish a message on rabbitmq.

Available settings:

- `RABBIT_USER` - Username
- `RABBIT_PASSWORD` - Password
- `RABBIT_HOST` - Hostname
- `RABBIT_PORT` - Port
- `RABBIT_VHOST` - Vhost
- `RABBIT_EXCHANGE` - Exchange name
- `RABBIT_QUEUE` - Queue

For defaults, look at `main.go`, lines 26-32.

## Usage

```sh
docker build -t mqp .
docker run mqp "my-message-content"
```
