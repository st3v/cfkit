package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/st3v/cfkit/env"
)

func main() {
	fmt.Println("ENVAPP", env.Addr())

	router := mux.NewRouter()

	router.HandleFunc("/app", appHandler)
	router.HandleFunc("/service/{name}", svcHandler)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(env.Addr(), nil))
}

func appHandler(rw http.ResponseWriter, req *http.Request) {
	app, err := env.Application()
	if err != nil {
		fmt.Fprint(rw, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(app)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(rw, string(out))
}

func svcHandler(rw http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	if name == "" {
		fmt.Fprint(rw, "Missing service name")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	svc, err := env.ServiceWithName(name)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	out, err := json.Marshal(svc)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(rw, string(out))
}
