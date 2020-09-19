# Nats-PubSub-Message
publish subscriber using nats and postgresql

How to run:
 1. edit properties in properties directory
 2. edit configPath in executor directory
 3. running executor(example: ". start-produce-service.sh")
 
How to test:
  1. running go script in test directory
 
Example:
  1.  running consume:
    - go run consume.go
        - response (status:success, error:), consume still waiting message by name queue
  2. running produce:
    - go run produce.go
        - response (status:success, error:), message will send to consume
