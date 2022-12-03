package web

import (
	"encoding/json"
	"net/http"
	"objects"
	"text/template"
)

type PageData struct {
	PageTitle string
	UrlPath   string
	Peoples   []objects.People
}

func ReadApiPeoples(w http.ResponseWriter, r *http.Request) []objects.People {
	resp, err := http.Get("http://localhost:5000/api/peoples")
	if err == nil {
		defer resp.Body.Close()

		var items []objects.People
		err := json.NewDecoder(resp.Body).Decode(&items)
		if err == nil {
			return items
		}
	}

	http.Error(w, err.Error(), http.StatusBadRequest)
	return nil
}

func Run() {
	tmpl := template.Must(template.ParseFiles("./web/layout.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := PageData{
			PageTitle: "Welcome",
			UrlPath:   r.URL.Path,
			Peoples:   ReadApiPeoples(w, r),
		}
		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
