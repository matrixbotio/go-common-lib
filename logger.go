package go_common_lib

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/matrixbotio/constants-lib"
	"log"
	"os"
	"strings"
)

type esLogger struct {
	esProto  string
	esHost   string
	esPort   string
	esIndex  string
	esClient *elasticsearch.Client
}

func InitLogger(sourceName string, logLevel string, esProto string, esHost string, esPort string,
	esIndex string) (*constants.Logger, error) {
	esLogger := esLogger{
		esProto: esProto,
		esHost:  esHost,
		esPort:  esPort,
		esIndex: esIndex,
	}

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	logger, err := constants.NewLogger(&esLogger, hostName, sourceName, logLevel), nil
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
	if l.esClient == nil {
		return
	}

	_, err := l.esClient.Index(
		l.esIndex,
		strings.NewReader(data),
		l.esClient.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Println("failed to send log to ES: " + err.Error())
	}
}

func (l *esLogger) initEs(logger *constants.Logger) error {
	if l.esProto == "" {
		logger.Log("ElasticSearch protocol is not passed, initialising logger without ElasticSearch")
		return nil
	}

	var err error
	esConfig := elasticsearch.Config{
		Addresses: []string{
			l.esProto + "://" + l.esHost + ":" + l.esPort,
		},
	}
	l.esClient, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		return err
	}
	return nil
}
