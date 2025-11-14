package common

import "encoding/json"

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	HttpPort string `json:"http_port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IndicatorSettings struct {
	Spread     string `json:"spread,omitempty"`
	Imbalances string `json:"imbalances,omitempty"`
	SecretKey  string `json:"secretkey,omitempty"`
}

type ConfigFile struct {
	Rabbitmq RabbitMQ                                `json:"rabbitmq,omitempty"`
	Sources  map[string]map[string]IndicatorSettings `json:"sources,omitempty"`
	Strategy json.RawMessage                         `json:"strategy,omitempty"`
}
