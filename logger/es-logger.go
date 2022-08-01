package logger

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"os"
	"strings"
)

func InitESLogger(sourceName string, logLevel string, esProto string, esHost string, esPort string,
	esIndex string) (*Logger, error) {
	esLogger := esLogger{
		proto: esProto,
		host:  esHost,
		port:  esPort,
		index: esIndex,
	}

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(&esLogger, hostName, sourceName, logLevel), nil
	if err != nil {
		return nil, err
	}

	err = esLogger.initEs(logger)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *esLogger) Send(data string) {
	if l.client == nil {
		return
	}

	_, err := l.client.Index(
		l.index,
		strings.NewReader(data),
		l.client.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Println("failed to send log to ES: " + err.Error())
	}
}

func (l *esLogger) initEs(logger *Logger) error {
	if l.proto == "" {
		logger.Log("ElasticSearch protocol is not passed, initialising logger without ElasticSearch")
		return nil
	}

	var err error
	esConfig := elasticsearch.Config{
		Addresses: []string{
			l.proto + "://" + l.host + ":" + l.port,
		},
	}
	l.client, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		return err
	}
	return nil
}
