package main

import (
	"fmt"
	. "github.com/CarlosRA97/goOpenServerReceiver/worker"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	fmt.Println(Configure())
	listDir, err := ioutil.ReadDir("images")
	Check(err)
	listDir = listDir[1:]
	for file := range listDir {
		ImageUploader(fmt.Sprintf("images/%s", listDir[file].Name()))
	}
	//ImageUploader("Go-aviator.png")
}

func main() {

	http.HandleFunc("/", Index)

	http.HandleFunc("/update", UpdateMapsOnServer)
	http.HandleFunc("/update-map", ChangeMap)
	http.HandleFunc("/server-state", State)

	port := Configure().Port

	fmt.Printf("Listening on Port: %s ...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil) // , "server.pem", "server.key"
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
