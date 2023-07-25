package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"seitaikei/utils"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var data = map[string][]byte{}
var endpoints = []string{}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	var jdata []interface{}
	for _, k := range endpoints {
		var tmp interface{}
		json.Unmarshal(data[k], &tmp)
		jdata = append(jdata, tmp)
	}
	//w.Write([]byte(tmp))

	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	tmpl.Execute(w, jdata)
}

func apiHandle(w http.ResponseWriter, r *http.Request) {
	endpnt := mux.Vars(r)["src"]

	if _, ok := data[endpnt]; !ok {
		log.Printf("Got data at api enpoint: %q", endpnt)
		endpoints = append(endpoints, endpnt)
	}
	data[endpnt], _ = io.ReadAll(r.Body)
}

func inputCtrl(done chan bool) {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadString('\n')
	done <- true
}

func main() {
	utils.Banner()
	log.Println("Starting Web Interface....")

	done := make(chan bool)

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandle)
	r.HandleFunc("/api/{src}", apiHandle).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Println(srv.ListenAndServe())
	}()
	go inputCtrl(done)
	<-done
	srv.Shutdown(context.Background())
	log.Println("Web Server Shutdown")
}
