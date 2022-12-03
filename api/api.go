package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"objects"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getIdFunc(obj any) any {
	return obj.(*objects.People).Id
}

var Peoples = NewRepository[objects.People](getIdFunc)

func seed() {
	Peoples.Add(&objects.People{Id: uuid.New(), Name: "Alex", Age: 47})
	Peoples.Add(&objects.People{Id: uuid.New(), Name: "Egor", Age: 35})
	Peoples.Add(&objects.People{Id: uuid.New(), Name: "Ivan", Age: 25})
	Peoples.Add(&objects.People{Id: uuid.New(), Name: "Kostya", Age: 15})
	Peoples.Add(&objects.People{Id: uuid.New(), Name: "Misha", Age: 20})
}

func WriteJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func registerRepositoryHandlers(rep *Repository[objects.People], route string) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		var item objects.People
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item.Id = uuid.New()
		ok := rep.Add(&item)
		WriteJsonResponse(w, ok)
	}).Methods("PUT")

	r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		var item objects.People
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok := rep.Update(&item)
		WriteJsonResponse(w, ok)
	}).Methods("POST")

	r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		WriteJsonResponse(w, rep.GetAll())
	}).Methods("GET")

	r.HandleFunc(route+"/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		if id, error := uuid.Parse(key); error == nil {
			if data, ok := rep.Get(id); ok {
				WriteJsonResponse(w, data)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}).Methods("GET")

	r.HandleFunc(route+"/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		if id, error := uuid.Parse(key); error == nil {
			if ok := rep.Remove(id); ok {
				WriteJsonResponse(w, ok)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}).Methods("DELETE")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	return r
}

func createServer(port int, handler http.Handler) *http.Server {
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: handler,
	}
	return &server
}

func Run(port int) {
	seed()

	handler := registerRepositoryHandlers(Peoples, "/api/peoples")
	server := createServer(port, handler)
	server.ListenAndServe()

}
