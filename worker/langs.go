package worker

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"golang.org/x/text/language"
)

type LoginPageStrings struct {
	Title               string `json:"title"`
	UsernameLabel       string `json:"username_label"`
	UsernamePlaceholder string `json:"username_placeholder"`
	PasswordLabel       string `json:"password_label"`
	LoginLabel          string `json:"login_label"`
}

type LoginPage struct {
	Login struct {
		Es LoginPageStrings `json:"es"`
		En LoginPageStrings `json:"en"`
	} `json:"login"`
}

func (l *LoginPage) ParseDataFrom(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(data, l)
}

func (l *LoginPage) GetLang(lang string) LoginPageStrings {
	var d LoginPageStrings
	switch lang {
	case language.Spanish.String():
		d = l.Login.Es
	case language.English.String():
		d = l.Login.En
	}
	return d
}
