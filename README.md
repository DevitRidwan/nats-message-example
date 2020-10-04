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
 5. running executor(example: ". start-endpoint-service.sh")

Doc Test => https://documenter.getpostman.com/view/8200038/TVRg6Umw
