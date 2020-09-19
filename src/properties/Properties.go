package properties

import (
	"Ridwan/Queue/src/dbms"
	"fmt"
)

type ServiceProperties struct {
	Nats        NatsProperties      `json:"nats" mapstructure:"nats"`
	Logging     LoggingProperties   `json:"logging" mapstructure:"logging"`
	NameLogging string              `json:"name_log" mapstructure:"name_log"`
	Database    *dbms.Properties    `json:"database"  mapstructure:"database"`
	Topic       NatsTopicProperties `json:"topic" mapstructure:"topic"`
}

type LoggingProperties struct {
	OutputPaths      []string `json:"output_paths" mapstructure:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths" mapstructure:"error_output_paths"`
}

type NatsTopicProperties struct {
	BaseURL            string                       `json:"base_url" mapstructure:"base_url"`
	EnvirontmentStatus string                       `json:"environtment_status" mapstructure:"environtment_status"`
	Channel            map[string]ChannelProperties `json:"channel" `
}

type ChannelProperties struct {
	Path        string `json:"path" mapstructure:"path"`
	QueueName   string `json:"queue_name" mapstructure:"queue_name"`
	DurableName string `json:"durable_name" mapstructure:"durable_name"`
}

type NatsProperties struct {
	Address   string `json:"address" mapstructure:"address" `
	ClusterId string `json:"cluster_id" mapstructure:"cluster_id"`
	ClientId  string `json:"client_id" mapstructure:"client_id"`
}

func (this *NatsTopicProperties) GetChannel(name string) string {
	return fmt.Sprintf("%s.%s.%s", this.BaseURL, this.EnvirontmentStatus, this.Channel[name].Path)
}

func (this *NatsTopicProperties) GetQueueName(name string) string {
	return this.Channel[name].QueueName
}

func (this *NatsTopicProperties) GetDurableName(name string) string {
	return this.Channel[name].DurableName
}

func (this *NatsTopicProperties) GetQueueSubscribeProperties(name string) (string, string) {
	return this.GetChannel(name), this.GetQueueName(name)
}
