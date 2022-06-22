package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water          int    `json:"water"`
	Wind           int    `json:"wind"`
	StatusCompiled string `json:"status_compiled"`
}

type DataStatus struct {
	Status Status `json:"status"`
}

var (
	staticHTMLFiles = "views/templates.html"
	jsonPath        = "views/status.json"
	newStatus       DataStatus
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		file, _ := ioutil.ReadFile(jsonPath)
		err := json.Unmarshal(file, &newStatus)
		if err != nil {
			fmt.Println("errMarshal", err)
			panic(err)
		}

		tpl, err := template.ParseFiles(staticHTMLFiles)
		if err != nil {
			fmt.Println("errTPL", err)
			return
		}

		context := Status{
			Water:          newStatus.Status.Water,
			Wind:           newStatus.Status.Wind,
			StatusCompiled: newStatus.Status.StatusCompiled,
		}

		go func() {
			for {
				newStatus.Status.Water = rand.Intn(100)
				newStatus.Status.Wind = rand.Intn(100)

				if newStatus.Status.Water < 5 || newStatus.Status.Wind < 6 {
					newStatus.Status.StatusCompiled = "AMAN"
				} else if (newStatus.Status.Water >= 6 && newStatus.Status.Water <= 8) || (newStatus.Status.Wind >= 7 && newStatus.Status.Wind <= 15) {
					newStatus.Status.StatusCompiled = "SIAGA"
				} else if newStatus.Status.Water > 8 || newStatus.Status.Wind > 15 {
					newStatus.Status.StatusCompiled = "BAHAYA"
				}

				jsonString, _ := json.Marshal(&newStatus)
				err := ioutil.WriteFile(jsonPath, jsonString, os.ModePerm)
				if err != nil {
					fmt.Println("errWriteToJson", err)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}()

		err = tpl.Execute(w, &context)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Listening on port 127.0.0.1:8088")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8088",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
