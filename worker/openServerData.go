package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type openServerData struct {
	MapChoosen    string   `json:"map_choosen"`
	MapsAvailable []string `json:"maps_available"`
	ServerIP      string   `json:"server_ip"`
	ServerState   bool     `json:"server_state"`
	Version       string   `json:"version"`
}

func (o *openServerData) getDataFrom(url string) {
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()

		htmlData, _ := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(htmlData, o)
	}
}

func (o *openServerData) putDataTo(url string) {
	data := o.MapsAvailable

	jn, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jn))
	body := strings.NewReader(string(jn))

	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}
