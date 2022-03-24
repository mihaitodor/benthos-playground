# Benthos consuming a Redpanda topic using TLS

## Create certs

```shell
> terraform init
> terraform apply -auto-approve
```

## Run docker-compose

```shell
> docker-compose up
```

## Produce some messages

```shell
> docker exec redpanda rpk topic create test_topic
> docker exec -it redpanda rpk topic produce test_topic
```

Type some message and hit return. The same messages will be emitted in the Benthos container logs.

You can also use the `rpk` tool directly from the host if you have it installed:

```shell
> rpk topic produce test_topic --tls-enabled --tls-cert ./certs/client-cert.pem --tls-key ./certs/client-privkey.pem --tls-truststore ./certs/ca-cert.pem
```
