package pages

import (
	"html/template"
	"net/http"
)

type About struct{}

func (a *About) Render(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("./website/about.html")
	err := t.Execute(w, a)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
