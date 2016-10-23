package worker

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
)

func Index(w http.ResponseWriter, r *http.Request) {

	img := []Image{}
	GetData(&img, "images")

	images := map[string]int{}

	for i := range img {
		images[img[i].ID] = i
	}

	image := img[images["images (2)"]]

	w.Write([]byte(fmt.Sprintf(`<html><img src="data:image/%s;base64,%s"></img></html>`, image.Format, image.Data)))

}

func State(w http.ResponseWriter, r *http.Request) {
	var info openServerData

	fullRequest := url.URL{
		Scheme:   https,
		Host:     firDatabaseURL,
		Path:     firebaseDB(config.Main),
		RawQuery: authQuery(apiKey),
	}

	info.getDataFrom(fullRequest.String())

	runServer := exec.Command("python", "OpenServer.py")
	closeServer := exec.Command("pkill", "java")

	if info.ServerState {
		err := runServer.Run()
		Check(err)
	} else {
		err := closeServer.Run()
		Check(err)
	}

}

func ChangeMap(w http.ResponseWriter, r *http.Request) {

	var info openServerData

	getRequest := url.URL{
		Scheme:   https,
		Host:     firDatabaseURL,
		Path:     firebaseDB(config.Main),
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

func UpdateMapsOnServer(w http.ResponseWriter, r *http.Request) {

	var info openServerData

	putRequest := url.URL{
		Scheme:   https,
		Host:     firDatabaseURL,
		Path:     firebaseDB(config.MapsAvailable),
		RawQuery: authQuery(apiKey),
	}
	fmt.Println(putRequest.String())

	info.MapsAvailable = getMapsFrom("./")

	info.putDataTo(putRequest.String())
}
