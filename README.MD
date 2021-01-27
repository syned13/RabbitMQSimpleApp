# RabbitMQSimple app

A simple app to test RabbitMQ message queuing.

### How to run

#### Step 1
- Create a .env file in the project's root and under the worker folder.
- Put a variable CLOUDAMQP_URL, in the form of *amqp://user:pass@host/vhost* pointing to the connection string of a rabbitMQ queue.

#### Step 2

Run the following commands:

```go run main.go```

```go run worker/main.go```

