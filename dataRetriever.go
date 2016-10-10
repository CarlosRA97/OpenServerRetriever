package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/CarlosRA97/go-toml"
)

var (
	config, err    = toml.LoadFile("config.toml")
	firDatabaseURL string
	apiKey         string
	port           string
)

const (
	https string = "https"
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

func init() {
	apiKey = config.Get("FIRrestApi.apiKey").(string)
	firDatabaseURL = config.Get("FIRrestApi.dbUrl").(string)
}

func main() {

	port = config.Get("server.port").(string)

	http.HandleFunc("/update", updateMapsOnServer)
	http.HandleFunc("/update-map", changeMap)
	http.HandleFunc("/server-state", state)

	fmt.Printf("Listening on Port: %s ...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil) // , "server.pem", "server.key"
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func state(w http.ResponseWriter, r *http.Request) {
	var info openServerData

	fullRequest := url.URL{
		Scheme:   https,
		Host:     firDatabaseURL,
		Path:     firebaseDB(config.Get("userDBPath.main").(string)),
		RawQuery: authQuery(apiKey),
	}

	info.getDataFrom(fullRequest.String())

	runServer := exec.Command("python", "OpenServer.py")
	closeServer := exec.Command("pkill", "java")

	if info.ServerState {
		err := runServer.Run()
		check(err)
	} else {
		err := closeServer.Run()
		check(err)
	}

}

func changeMap(w http.ResponseWriter, r *http.Request) {

	var info openServerData

	getRequest := url.URL{
		Scheme:   https,
		Host:     firDatabaseURL,
		Path:     firebaseDB(config.Get("userDBPath.main").(string)),
		RawQuery: authQuery(apiKey),
	}

	fmt.Println(getRequest.String())

	info.getDataFrom(getRequest.String())

	mapName := info.MapChoosen

	command := exec.Command("python", "pMod.py", mapName)
	err := command.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func updateMapsOnServer(w http.ResponseWriter, r *http.Request) {

	var info openServerData

	putRequest := url.URL{
		Scheme:   https,
		Host:     firDatabaseURL,
		Path:     firebaseDB(config.Get("userDBPath.mapsAvailable").(string)),
		RawQuery: authQuery(apiKey),
	}
	fmt.Println(putRequest.String())

	info.MapsAvailable = getMapsFrom("./")

	info.putDataTo(putRequest.String())
}

func firebaseDB(path string) string {
	return fmt.Sprintf("%s.json", path)
}

func authQuery(key string) string {
	return fmt.Sprintf("auth=%s", key)
}

func getMapsFrom(path string) []string {
	var maps []string
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		nonMapsFolders := f.Name() != ".DS_Store" && f.Name() != ".idea" && f.Name() != "logs" && f.Name() != "mapas-originales"
		if nonMapsFolders {
			if f.IsDir() {
				maps = append(maps, f.Name())
			}
		}
	}
	return maps
}
