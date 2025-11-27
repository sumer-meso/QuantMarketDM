package common

import "encoding/json"

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	HttpPort string `json:"http_port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountSettings struct {
	AccountID string `json:"account_id,omitempty"`
	ApiKey    string `json:"api_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
}

type IndicatorSettings struct {
	Spread     string `json:"spread,omitempty"`
	Imbalances string `json:"imbalances,omitempty"`
}

type SourceSettings struct {
	Accounts   []AccountSettings            `json:"accounts,omitempty"`
	Indicators map[string]IndicatorSettings `json:"indicators,omitempty"`
}

type ConfigFile struct {
	Rabbitmq RabbitMQ                  `json:"rabbitmq"`
	Sources  map[string]SourceSettings `json:"sources"`
	Strategy json.RawMessage           `json:"strategy,omitempty"`
}
