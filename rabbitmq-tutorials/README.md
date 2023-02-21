# RabbitMQ Demo

https://www.rabbitmq.com/tutorials/

https://github.com/rabbitmq/rabbitmq-tutorials

## Run RabbitMQ on Docker
```
# latest RabbitMQ 3.11
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management
```