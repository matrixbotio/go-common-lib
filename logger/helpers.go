package logger

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func getJSON(url string, storage interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln("Unable to download errors JSON file: " + url)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicln("Exception while closing errors JSON body: " + err.Error())
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("Exception while reading errors JSON body: " + err.Error())
		return
	}
	err = json.Unmarshal(body, storage)
	if err != nil {
		log.Panicln("Exception while unmarshalling errors JSON body: " + err.Error())
		return
	}
}
