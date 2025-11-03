package common

import "encoding/json"

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IndicatorSettings struct {
	Spread     string `json:"spread,omitempty"`
	Imbalances string `json:"imbalances,omitempty"`
	SecretKey  string `json:"secretkey,omitempty"`
}

type ConfigFile struct {
	RabbitmqTest RabbitMQ                                `json:"rabbitmqTest"`
	RabbitmqProd RabbitMQ                                `json:"rabbitmqProd"`
	Sources      map[string]map[string]IndicatorSettings `json:"sources"`
	Strategy     json.RawMessage                         `json:"strategy"`
}
