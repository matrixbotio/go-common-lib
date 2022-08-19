package logger

import "github.com/elastic/go-elasticsearch/v7"

type logDevice interface {
	Send(data string)
}

type esLogger struct {
	proto  string
	host   string
	port   string
	index  string
	client *elasticsearch.Client
}

type sendMessageFormat struct {
	Source    string      `json:"source"`
	Host      string      `json:"host"`
	Timestamp string      `json:"timestamp"`
	Level     int         `json:"level"`
	Message   string      `json:"message"`
	Code      interface{} `json:"code,omitempty"`
	Stack     interface{} `json:"stack,omitempty"`
}

type logLevelDesc struct {
	Stderr bool
	Format string
	Level  int
}

type Logger struct {
	Dev           logDevice
	Host          string
	Source        string
	LowestLevel   int
	PartLogLevels map[string]int
}

type logConfiguration struct {
	LogLevels   map[string]*logLevelDesc
	DTFormat    string
	DTFormatLen int
}
