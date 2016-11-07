package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"io/ioutil"
	"golang.org/x/crypto/bcrypt"
	"io"
)

const (
	usernameField = "username"
	passwordField = "password"
	GET           = "GET"
	POST          = "POST"
	loginHTML     = "login.html"
	indexHTML     = "index.html"
)

var (
	tpl  = LoadTemplates()
	user = make(map[string][]byte)
	//username string
	//password string
)

func Index(w http.ResponseWriter, r *http.Request) {

	lang := strings.Split(strings.Split(r.Header.Get("Accept-Language"), ",")[0], "-")[0]

	l := LoginPage{}
	l.ParseDataFrom("./login_lang.json")

	loginPageStrings := l.GetLang(lang)

	switch r.Method {
	case GET:
		tpl.ExecuteTemplate(w, loginHTML, loginPageStrings)
	case POST:
		if ok := verifyLogin(r.FormValue(usernameField), r.FormValue(passwordField)); ok {
			tpl.ExecuteTemplate(w, indexHTML, nil)
		} else {
			tpl.ExecuteTemplate(w, loginHTML, loginPageStrings)
		}
	}
}

func verifyLogin(username string, password string) bool {
	var ok bool

	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(data, &user)

	if bcrypt.CompareHashAndPassword(user[username], []byte(password)) == nil {
		ok = true
	}
	return ok
}

func Register(w http.ResponseWriter, r *http.Request) {
	username, password, auth := r.URL.Query().Get("username"), r.URL.Query().Get("password"), r.URL.Query().Get("auth")

	if auth == Configure().Server.ApiKey {
		hahedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			log.Fatalln(err)
		}

		user[username] = hahedPassword

		data, err := json.Marshal(user)

		err = WriteFile("users.json", strings.NewReader(string(data)))
		if err != nil {
			log.Fatalln(err)
		}
		w.Write([]byte("Registered Correctly"))
	} else {
		w.Write([]byte("Failed to register"))
	}

}

func WriteFile(filename string, r io.Reader) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
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

	changeMapCommand(mapName)
}

func changeMapCommand(mapName string) {
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
