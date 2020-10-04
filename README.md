
# Nats-PubSub-Message
publish subscriber using nats and postgresql

Libraries:
- github.com/jinzhu/gorm
- github.com/nats-io/nats.go
- github.com/nats-io/stan.go
- go.uber.org/zap
- github.com/spf13/viper


How to run:
 1. import db_queue.sql into your database
 2. download libraries go
 3. edit properties in properties directory
 4. edit configPath in executor directory
 5. running executor(example: ". start-produce-service.sh")
 
How to test:
  1. running go script in test directory
 
Example:
  1.  running consume:
    - go run consume.go
        - response (status:success, error:), consume still waiting message by name queue
  2. running produce:
    - go run produce.go
        - response (status:success, error:), message will send to consume
