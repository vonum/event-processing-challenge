# Event processing task
## Setup
1. `go mod download`
2. `make run`
3. `make migrate`

## Running components
1. Run publisher -> `go run internal/cmd/publisher/main.go`
2. Run subscribers -> `go run internal/cmd/subscriber/main.go`

## Event processing
RabbitMQ was chosen as a message broker for the following reasons:
1. Simplicity
2. Good defaults for a message queue implementation
3. No previous experience with RabbitMQ

Messages are encoded using protobuf to improve throughput.

Messages are published to a queue and subscribers take messages in round robin
order.

Message acknowledgment is done manually to avoid messages being lost due to
failing tasks. If a connection is lost or a timeout is exceeded, RabbitMQ will
dispatch the message to another subscriber.

In case processing fails for a message, acknowledgment is done on the queue and
message passed to a dead letter queue.

Messages and queues are made durable so they are kept on disk even after
RabbitMQ restart.

## Caching strategy
Redis was chosen for caching. All exchange rates are stored together for
simplicity.

To prevent race conditions and spamming the external api. Redis is also used as
a distributed lock.

This means when a worker doesn't find values in cache, it will try to acquire a
lock to consume the api.

If it succeeds, the worker will consume the api and write data to cache.

If it doesn't, it will continue looping with a timeout:
1. read from cache
2. try to acquire a lock

If a specified timeout is reached, the worker will fail and write the message to
dead letter queue.

Quotes will expire with a specified TTL.
Locks will also expire to prevent deadlocks due to workers holding the lock
failing.

## Api
Api has a lot of room for improvement as it was not the focus. It is done just
to be done.
You can consume the api with:
1. `curl localhost:3000/materialize`
2. `curl localhost:3000/materialize?window_size=10`

## Further improvements
1. RabbitMQ configuration could be further explored to tune the setup.
2. Much better error handling -> currently only a few errors are handled while
   the rest are not due to time issues.
3. Tests are also missing due to time issues.
4. All components should be containerized and deployed with a container
   orchestration service like kubernetes.
5. The api is very minimalistic and should be refactored
6. The api should have a go routine that would free memory for events per second
   to avoid consuming all process memory
7. Adding context to logs, like worker id, event id etc. This would make log
   aggregation easier. Current logging setup is made for easier monitoring of
   components locally.
