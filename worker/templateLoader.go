package worker

import "html/template"

func LoadTemplates() *template.Template {
	return template.Must(template.ParseGlob("./templates/*.html"))
}
